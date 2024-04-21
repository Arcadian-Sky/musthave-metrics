package controller

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

			err = c.sendPack(metrics, pollCount)
			if err != nil {
				fmt.Println("Error sending metrics:", err)
			}

			pollCount = pollCount + 1
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

func (c *CollectAndSendMetricsService) makePack(metrics map[string]interface{}, pollCount int) []interface{} {
	forSend := make([]interface{}, 0, len(metrics))
	for metricType, value := range metrics {
		mValue := value.(float64)
		forSend = append(forSend, models.Metrics{
			ID:    metricType,
			MType: "gauge",
			Value: &mValue,
		})
	}
	mValue := int64(pollCount)
	forSend = append(forSend, models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &mValue,
	})

	return forSend
}

// Отправляем метрики
func (c *CollectAndSendMetricsService) send(metrics map[string]interface{}, pollCount int) error {
	var forSend = c.makePack(metrics, pollCount)
	for _, metric := range forSend {
		err := c.sendMetricJSON(metric, "/update")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CollectAndSendMetricsService) sendPack(metrics map[string]interface{}, pollCount int) error {
	var forSend = c.makePack(metrics, pollCount)
	err := c.sendMetricJSON(forSend, "/updates")
	if err != nil {
		return err
	}

	return nil
}

// Отправляем запрос на сервер
func (c *CollectAndSendMetricsService) sendMetricJSON(m any, method string) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error marshaling metrics:", err)
		return err
	}
	fmt.Printf("m: %v\n", string(jsonData))

	// Формируем адрес запроса
	url := fmt.Sprintf("%s"+method, c.config.GetServerAddress())

	// Создание HTTP-запроса POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	hashKey := c.config.GetHash()
	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(jsonData)
		dst := h.Sum(nil)
		// fmt.Printf("dst: %v\n", hex.EncodeToString(dst))
		req.Header.Set("HashSHA256", hex.EncodeToString(dst))
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//
	// Отправляем запрос на сервер
	// resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	fmt.Printf("Metrics did not sent: \n")
	// 	return err
	// }
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
	// fmt.Printf("Metric sent: %s\n", mName)
	defer resp.Body.Close()

	return nil
}
