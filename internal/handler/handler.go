package handler

import (
	"fmt"
	"net/http"

	"github.com/Arcadian-Sky/musthave-metrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

func MetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if storage.Storage == nil {
		storage.Storage = storage.NewMemStorage()
	}
	currentMetrics := storage.Storage.GetMetrics()
	for name, value := range currentMetrics {
		fmt.Fprintf(w, "%d: %v\n", name, value)
	}
}

func UpdateMetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if storage.Storage == nil {
		storage.Storage = storage.NewMemStorage()
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	// fmt.Fprintf(w, "storage %s\n", metricType)
	// fmt.Fprintf(w, "storage %s\n", metricName)
	// fmt.Fprintf(w, "storage %s\n", metricValue)
	fmt.Println("nonEmptyParts", metricType, metricName, metricValue, "|")
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

func GetMetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if storage.Storage == nil {
		storage.Storage = storage.NewMemStorage()
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	// fmt.Fprintf(w, "storage %s\n", metricType)
	// fmt.Fprintf(w, "storage %s\n", metricName)
	fmt.Println("nonEmptyParts", metricType, metricName, "|")
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

	// Выводим данные
	metricTypeID, err := storage.GetMetricTypeByCode(metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentMetrics := storage.Storage.GetMetric(metricTypeID)
	fmt.Printf("metricName: %v\n", metricName)
	if metricName != "" {
		fmt.Printf("currentMetrics[metricName]: %v\n", currentMetrics[metricName])
		if currentMetrics[metricName] != nil {
			_, err = w.Write([]byte(fmt.Sprintf("%v", currentMetrics[metricName])))
			if err != nil {
				http.Error(w, "w.Write Error: "+err.Error(), http.StatusNotFound)
				//return err
			}
		} else {
			http.Error(w, "Metric value not provided", http.StatusNotFound)
		}
	} else {
		for name, value := range currentMetrics {
			fmt.Fprintf(w, "%s: %v\n", name, value)
		}
	}

	w.WriteHeader(http.StatusOK)
}
