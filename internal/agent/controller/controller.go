package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
)

var (
	metrics map[string]interface{}
)

func CollectAndSendMetrics() {
	var pollCount int

	fmt.Println(flags.GetPollInterval(), flags.GetReportInterval())
	// Отправляем метрики на сервер
	go func() {
		fmt.Println("send")
		for {
			fmt.Println("n=", pollCount)
			err := send(metrics, pollCount)
			if err != nil {
				fmt.Println("Error sending metrics:", err)
			}
			pollCount = pollCount + 1
			fmt.Println("send2")
			time.Sleep(flags.GetReportInterval())
		}

	}()

	//Собираем метрики
	go func() {
		for {
			fmt.Println("updateMterics")
			collectMetrics(pollCount)
			time.Sleep(flags.GetPollInterval())
		}
	}()

	// go sendMetrics()
	// go updateMterics()

	select {}
}

// Собираем метрики
func collectMetrics(pollCount int) map[string]interface{} {
	metrics = make(map[string]interface{})

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
	metrics["PollCount"] = pollCount
	metrics["RandomValue"] = rand.Float64() // Произвольное значение

	return metrics
}

// Отправляем метрики на сервер
func send(metrics map[string]interface{}, pollCount int) error {
	for metricType, value := range metrics {
		err := sendMetricValue("gauge", metricType, value)
		if err != nil {
			return err
		}
	}
	err := sendMetricValue("count", "PollCount", pollCount)
	if err != nil {
		return err
	}

	return nil
}

func sendMetricValue(mType string, mName string, mValue interface{}) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Формируем адрес запроса
	url := fmt.Sprintf("%s/update/"+mType+"/%s/%v", flags.GetServerAddress(), mName, mValue)

	// Отправляем запрос на сервер
	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Printf("Metric did not sent: %s\n", mName)
		return err
	}

	// Печатаем результат отправки (для демонстрации, лучше использовать логгер)
	fmt.Printf("Metric sent: %s\n", mName)
	defer resp.Body.Close()

	return nil
}
