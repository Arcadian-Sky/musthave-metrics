package server

import (
	"context"
	"log"
	"net/http"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	appgrpc "github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/grpc"
	apphttp "github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/http"
	pb "github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/protometrics"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/router"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
