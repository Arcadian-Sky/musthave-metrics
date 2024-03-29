package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/repository"
)

type CollectAndSendMetricsService struct {
	config flags.Config
}

func NewCollectAndSendMetricsService(config flags.Config) *CollectAndSendMetricsService {
	return &CollectAndSendMetricsService{config}
}

func (c *CollectAndSendMetricsService) Run() {
	var pollCount int
	metricsRepo := repository.NewInMemoryMetricsRepository()

	// Отправляем метрики на сервер
	go func() {
		// fmt.Println("send")
		for {
			metrics, err := metricsRepo.GetMetrics()
			if err != nil {
				fmt.Println("Error collecting metrics:", err)
				return
			}
			err = c.send(metrics, pollCount)
			if err != nil {
				fmt.Println("Error sending metrics:", err)
			}
			pollCount = pollCount + 1
			fmt.Println("send2")
			time.Sleep(c.config.GetPollInterval())
		}

	}()

	//Собираем метрики
	go func() {
		for {
			// fmt.Println("updateMterics")
			_, err := metricsRepo.GetMetrics()
			if err != nil {
				// fmt.Println("Error collecting metrics:", err)
				return
			}
			time.Sleep(c.config.GetPollInterval())
		}
	}()

	select {}
}

// Отправляем метрики
func (c *CollectAndSendMetricsService) send(metrics map[string]interface{}, pollCount int) error {
	for metricType, value := range metrics {
		mValue := value.(float64)
		metric := models.Metrics{
			ID:    metricType,
			MType: "gauge",
			Value: &mValue,
		}
		err := c.sendMetricJSONValue(metric)
		if err != nil {
			return err
		}
	}
	mValue := int64(pollCount)
	metric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &mValue,
	}
	err := c.sendMetricJSONValue(metric)
	if err != nil {
		return err
	}

	return nil
}

// Отправляем запрос на сервер
func (c *CollectAndSendMetricsService) sendMetricJSONValue(m models.Metrics) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error marshaling metrics:", err)
		return err
	}

	// Формируем адрес запроса
	url := fmt.Sprintf("%s/update", c.config.GetServerAddress())

	// Отправляем запрос на сервер
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Metrics did not sent: %s\n", m.ID)
		return err
	}
	fmt.Printf("Metric sent: %s\n", m.ID)
	defer resp.Body.Close()

	return nil
}

func (c *CollectAndSendMetricsService) sendMetricValue(mType string, mName string, mValue interface{}) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Формируем адрес запроса
	url := fmt.Sprintf("%s/update/"+mType+"/%s/%v", c.config.GetServerAddress(), mName, mValue)

	// Отправляем запрос на сервер
	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		fmt.Printf("Metric did not sent: %s\n", mName)
		return err
	}

	// Печатаем результат отправки (для демонстрации, лучше использовать логгер)
	fmt.Printf("Metric sent: %s\n", mName)
	defer resp.Body.Close()

	return nil
}
