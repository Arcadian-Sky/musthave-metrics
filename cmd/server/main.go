package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/server"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
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

func main() {
	// Обработчик сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	parsed := flags.Parse()
	storeMetrics := storage.NewMemStorage(storage.Config{
		Interval:        parsed.StoreInterval,
		FileStoragePath: parsed.FileStorage,
		Restore:         parsed.RestoreMetrics,
	})
	vhandler := handler.NewHandler(storeMetrics)

	go func() {
		err := storeMetrics.LoadMetrics()
		// logger := packmiddleware.GetLogger()
		// logger.Info("Server started")

		if err != nil {
			log.Fatal(err.Error())
		}
		go storeMetrics.SaveMetrics()
		fmt.Println("Server started")
		log.Fatal(http.ListenAndServe(parsed.Endpoint, server.InitRouter(*vhandler)))
	}()

	<-stop
	storeMetrics.SaveToFile()
	fmt.Println("Server stopped")
}
