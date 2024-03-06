package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Arcadian-Sky/musthave-metrics/internal/storage"
)

// metricsHandler обрабатывает HTTP запросы для получения текущих метрик
func MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if storage.Storage == nil {
			storage.Storage = storage.NewMemStorage()
		}
		currentMetrics := storage.Storage.GetMetrics()
		for name, value := range currentMetrics {
			fmt.Fprintf(w, "%d: %v\n", name, value)
		}
	}
}

// updateMetricsHandler обрабатывает HTTP запросы для обновления метрик
func UpdateMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if storage.Storage == nil {
			storage.Storage = storage.NewMemStorage()
		}
		// fmt.Fprintf(w, "storage %s\n", memStorage)

		path := r.URL.Path
		// Разбиваем путь на части по слэшам
		parts := strings.Split(path, "/")

		var nonEmptyParts []string
		for _, part := range parts {
			if part != "" {
				nonEmptyParts = append(nonEmptyParts, part)
			}
		}
		// fmt.Println("nonEmptyParts \n", nonEmptyParts[0], nonEmptyParts[1], len(nonEmptyParts))
		length := len(nonEmptyParts)
		metricType := ""
		metricName := ""
		metricValue := ""
		if 2 <= length {
			metricType = nonEmptyParts[1]
		}
		if 3 <= length {
			metricName = nonEmptyParts[2]
		}
		if 4 <= length {
			metricValue = nonEmptyParts[3]
		}
		//Проверяем передачу типа
		if metricType == "" {
			http.Error(w, "Metric type not provided", http.StatusNotFound)
			return
		}
		//Проверяем передачу имени
		if metricName == "" {
			http.Error(w, "Metric name not provided", http.StatusNotFound)
			return
		}

		// Обновляем метрику
		err := storage.Storage.UpdateMetric(metricType, metricName, metricValue)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Выводим данные
		currentMetrics := storage.Storage.GetMetrics()
		for name, value := range currentMetrics {
			fmt.Fprintf(w, "%d: %v\n", name, value)
		}

		w.WriteHeader(http.StatusOK)
	}
}
