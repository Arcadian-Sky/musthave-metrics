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
		if storage.Storage == nil {
			storage.Storage = storage.NewMemStorage()
		}
		currentMetrics := storage.Storage.GetMetrics()
		for name, value := range currentMetrics {
			fmt.Fprintf(w, "%s: %v\n", name, value)
		}
	}
}

// updateMetricsHandler обрабатывает HTTP запросы для обновления метрик
func UpdateMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if storage.Storage == nil {
			storage.Storage = storage.NewMemStorage()
		}
		// fmt.Fprintf(w, "storage %s\n", memStorage)

		path := r.URL.Path
		// Разбиваем путь на части по слэшам
		parts := strings.Split(path, "/")

		// fmt.Fprintf(w, "parts %s\n", parts)

		metricType := parts[2]
		metricName := parts[3]
		metricValue := parts[4]
		// fmt.Fprintf(w, "metrics %s %s %s", metricType, metricName, metricValue)

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
			fmt.Fprintf(w, "%s: %v\n", name, value)
		}

		w.WriteHeader(http.StatusOK)
	}
}
