package repository

import (
	"log"
	"math/rand"
	"runtime"
)

// MetricsRepository определяет методы для работы с метриками.
type MetricsRepository interface {
	GetMetrics() (map[string]interface{}, error)
}

type InMemoryMetricsRepository struct {
	metrics map[string]interface{}
}

func NewInMemoryMetricsRepository() *InMemoryMetricsRepository {
	return &InMemoryMetricsRepository{
		metrics: make(map[string]interface{}),
	}
}

func (r *InMemoryMetricsRepository) GetMetrics() (map[string]interface{}, error) {

	metrics := make(map[string]interface{})

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Собираем метрики из пакета runtime
	metrics["Alloc"] = float64(memStats.Alloc)
	metrics["BuckHashSys"] = float64(memStats.BuckHashSys)
	metrics["Frees"] = float64(memStats.Frees)
	metrics["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	metrics["GCSys"] = float64(memStats.GCSys)
	metrics["HeapAlloc"] = float64(memStats.HeapAlloc)
	metrics["HeapIdle"] = float64(memStats.HeapIdle)
	metrics["HeapIdle"] = float64(memStats.HeapIdle)
	metrics["HeapInuse"] = float64(memStats.HeapInuse)
	metrics["HeapObjects"] = float64(memStats.HeapObjects)
	metrics["HeapReleased"] = float64(memStats.HeapReleased)
	metrics["HeapSys"] = float64(memStats.HeapSys)
	metrics["LastGC"] = float64(memStats.LastGC)
	metrics["Lookups"] = float64(memStats.Lookups)
	metrics["MCacheInuse"] = float64(memStats.MCacheInuse)
	metrics["MCacheSys"] = float64(memStats.MCacheSys)
	metrics["MSpanInuse"] = float64(memStats.MSpanInuse)
	metrics["MSpanSys"] = float64(memStats.MSpanSys)
	metrics["Mallocs"] = float64(memStats.Mallocs)
	metrics["NextGC"] = float64(memStats.NextGC)
	metrics["NumForcedGC"] = float64(memStats.NumForcedGC)
	metrics["NumGC"] = float64(memStats.NumGC)
	metrics["OtherSys"] = float64(memStats.OtherSys)
	metrics["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	metrics["StackInuse"] = float64(memStats.StackInuse)
	metrics["StackSys"] = float64(memStats.StackSys)
	metrics["Sys"] = float64(memStats.Sys)
	metrics["TotalAlloc"] = float64(memStats.TotalAlloc)
	// Добавляем дополнительные метрики
	metrics["RandomValue"] = rand.Float64() // Произвольное значение

	// v, _ := mem.VirtualMemory()

	// metrics["TotalMemory"] = v.Total
	// metrics["FreeMemory"] = v.Free
	// v, _ := cpu.Counts(false)
	// metrics["CPUutilization1"] = v.cpu_count

	err := r.SaveMetrics(metrics)
	if err != nil {
		log.Println("Error saving metrics:", err)
	}
	return metrics, nil
}

// SaveMetrics сохраняет метрики в хранилище.
func (r *InMemoryMetricsRepository) SaveMetrics(metrics map[string]interface{}) error {
	r.metrics = metrics
	return nil
}
