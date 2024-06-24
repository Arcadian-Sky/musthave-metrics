// Модуль server обеспечивает хранение метрик
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

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
	// Создаем канал для обработки сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	parsed := flags.Parse()

	db, err := sql.Open("pgx", parsed.DBSettings)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.Migrations)

	// NewMemStorage создает новый экземпляр хранилищв
	// Создаем хранилище
	storeMetrics := func(cnf *flags.InitedFlags, db *sql.DB) storage.MetricsStorage {
		if cnf.StorageType == "postgres" {
			return postgres.NewPostgresStorage(db)
		}
		// mementoStore = storeMetrics
		return inmemory.NewMemStorage()
	}(parsed, db)

	memStore, memStoreOk := storeMetrics.(config.MementoStorage)
	if memStoreOk {
		// Инициализируем конфигурацию
		err = config.InitConfig(memStore, parsed)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	fmt.Println("cnf.StorageType:", parsed.StorageType)
	fmt.Println("memStoreOk:", memStoreOk)

	//Инициируем хендлеры
	vhandler := handler.NewHandler(storeMetrics, parsed)
	go func() {
		log.Println("Starting pprof server...")
		log.Fatal(http.ListenAndServe(parsed.Endpoint, server.InitRouter(*vhandler)))
	}()

	<-stop
	if memStoreOk {
		config.SaveMetricsToFile(memStore, parsed.FileStorage)
	}
	fmt.Println("Server stopped gracefully")
}
