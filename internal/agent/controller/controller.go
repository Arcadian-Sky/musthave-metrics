package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/repository"
)

func CollectAndSendMetrics() {
	var pollCount int
	metricsRepo := repository.NewInMemoryMetricsRepository()

	fmt.Println(flags.GetPollInterval(), flags.GetReportInterval())
	// Отправляем метрики на сервер
	go func() {
		fmt.Println("send")
		for {
			metrics, err := metricsRepo.GetMetrics()
			if err != nil {
				fmt.Println("Error collecting metrics:", err)
				return
			}
			err = send(metrics, pollCount)
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
			_, err := metricsRepo.GetMetrics()
			if err != nil {
				fmt.Println("Error collecting metrics:", err)
				return
			}
			time.Sleep(flags.GetPollInterval())
		}
	}()

	// go sendMetrics()
	// go updateMterics()

	select {}
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
