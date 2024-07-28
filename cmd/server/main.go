package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/server"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/config"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/inmemory"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/postgres"
	"github.com/Arcadian-Sky/musthave-metrics/migrations"
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

	err, memStore, memStoreOk := InitializeConfig(storeMetrics, parsed)
	if err != nil {
		log.Fatal(err.Error())
	}

	httpserver := InitializeHTTPServer(parsed, storeMetrics)

	go func() {
		log.Println("Starting server...")
		fmt.Printf("Build version: %s\n", buildVersion)
		fmt.Printf("Build date: %s\n", buildDate)
		fmt.Printf("Build commit: %s\n", buildCommit)
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

func InitializeConfig(storeMetrics storage.MetricsStorage, parsed *flags.InitedFlags) (error, config.MementoStorage, bool) {
	memStore, memStoreOk := storeMetrics.(config.MementoStorage)
	if storeMetrics == nil {
		return errors.New("storemetrics is nil"), nil, false
	}
	if memStoreOk {
		err := config.InitConfig(memStore, parsed)
		if err != nil {
			return err, nil, false
		}
	}

	log.Println("cnf.StorageType:", parsed.StorageType)
	log.Println("memStoreOk:", memStoreOk)

	return nil, memStore, memStoreOk
}

// Инициируем хендлеры
func InitializeHTTPServer(parsed *flags.InitedFlags, storeMetrics storage.MetricsStorage) *http.Server {
	vhandler := handler.NewHandler(storeMetrics, parsed)
	httpserver := &http.Server{
		Addr:    parsed.Endpoint,
		Handler: server.InitRouter(*vhandler, *parsed),
	}
	return httpserver
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
