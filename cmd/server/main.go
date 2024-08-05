package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pressly/goose/v3"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	appgrpc "github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/grpc"
	apphttp "github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/http"
	pb "github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/protometrics"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/router"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/config"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/inmemory"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/postgres"
	"github.com/Arcadian-Sky/musthave-metrics/migrations"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/grpclog"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// Пример запроса к серверу:
// POST /update/counter/someMetric/527 HTTP/1.1
// Host: localhost:8080
// Content-Length: 0
// Content-Type: text/plain

// Пример ответа от сервера:
// HTTP/1.1 200 OK
// Date: Tue, 21 Feb 2023 02:51:35 GMT
// Content-Length: 11
// Content-Type: text/plain; charset=utf-8
//
//	m := map[storage.MetricType]map[string]interface{}{
//		storage.Gauge: {
//			"metric1": 10.5,
//			"metric2": 20.7,
//		},
//		storage.Counter: {
//			"metric1": 100,
//			"metric2": 200,
//		},
//	}
//
// storeMetrics.SetMetrics(m)
func main() {
	stop := InitSignalHandler()

	parsed := flags.Parse()

	db, err := OpenDatabase(parsed.DBSettings)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.Migrations)

	storeMetrics := CreateMetricsStorage(parsed, db)

	memStore, memStoreOk, err := InitializeConfig(storeMetrics, parsed)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	if parsed.TcpEnable == true {
		grpcserver := InitializeGRPCServer(parsed, storeMetrics)
		go func() {
			log.Println("Starting GRPC server...")

			// определяем порт для сервера
			listen, err := net.Listen("tcp", parsed.TEndpoint)
			if err != nil {
				log.Fatalf("GRPC failed to listen: %v", err)
			}
			// получаем запрос gRPC
			if err := grpcserver.Serve(listen); err != nil {
				log.Fatalf("GRPC failed to serve: %v", err)
			}

		}()
	}

	httpserver := InitializeHTTPServer(parsed, storeMetrics)
	if parsed.TcpEnable == true {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		httpserver = InitializeTCP2HTTPServer(parsed, storeMetrics, ctx)
	}

	go func() {
		log.Println("Starting HTTP server...")
		if err := httpserver.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-stop

	// Handle graceful shutdown
	GracefulShutdown(httpserver, memStore, memStoreOk, parsed)
}

func InitSignalHandler() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}

func OpenDatabase(dbSettings string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbSettings)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateMetricsStorage(cnf *flags.InitedFlags, db *sql.DB) storage.MetricsStorage {
	// NewMemStorage создает новый экземпляр хранилищв
	// Создаем хранилище
	if cnf.StorageType == "postgres" {
		if db == nil {
			log.Println("CreateMetricsStorage: db is nil")
		}
		return postgres.NewPostgresStorage(db)
	}
	// mementoStore = storeMetrics
	return inmemory.NewMemStorage()
}

func InitializeConfig(storeMetrics storage.MetricsStorage, parsed *flags.InitedFlags) (config.MementoStorage, bool, error) {
	memStore, memStoreOk := storeMetrics.(config.MementoStorage)
	if storeMetrics == nil {
		return nil, false, errors.New("storemetrics is nil")
	}
	if memStoreOk {
		err := config.InitConfig(memStore, parsed)
		if err != nil {
			return nil, false, err
		}
	}

	log.Println("cnf.StorageType:", parsed.StorageType)
	log.Println("memStoreOk:", memStoreOk)

	return memStore, memStoreOk, nil
}

// Инициируем хендлеры
func InitializeTCP2HTTPServer(parsed *flags.InitedFlags, storeMetrics storage.MetricsStorage, ctx context.Context) *http.Server {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterMetricsServiceHandlerFromEndpoint(ctx, mux, parsed.TEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC Gateway handler: %v", err)
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	httpserver := &http.Server{
		Addr:    parsed.Endpoint,
		Handler: mux,
	}
	return httpserver
}

func InitializeHTTPServer(parsed *flags.InitedFlags, storeMetrics storage.MetricsStorage) *http.Server {
	vhandler := apphttp.NewHandler(storeMetrics, parsed)
	httpserver := &http.Server{
		Addr:    parsed.Endpoint,
		Handler: router.InitRouter(*vhandler, *parsed),
	}
	return httpserver
}

func InitializeGRPCServer(parsed *flags.InitedFlags, storeMetrics storage.MetricsStorage) *grpc.Server {
	metricsServer := appgrpc.NewServer(storeMetrics, parsed)
	// создаём gRPC-сервер без зарегистрированной службы
	grpcServer := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServiceServer(grpcServer, metricsServer)

	// ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	// Создайте HTTP/gRPC Gateway сервер
	// mux := runtime.NewServeMux()
	// err := pb.RegisterMetricsServiceHandlerServer(ctx, mux, metricsServer)
	return grpcServer
}

func GracefulShutdown(httpserver *http.Server, memStore config.MementoStorage, memStoreOk bool, parsed *flags.InitedFlags) {
	// Timeout for active connections to close
	shutdownTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpserver.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	if memStoreOk {
		config.SaveMetricsToFile(memStore, parsed.FileStorage)
	}

	fmt.Println("Server stopped gracefully")
}
