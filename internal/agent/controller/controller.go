package controller

import (
	"fmt"
	"time"

	senderPack "github.com/Arcadian-Sky/musthave-metrics/internal/agent/controller/sender"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/repository"
)

type CollectAndSendMetricsService struct {
	config flags.Config
	sender *senderPack.Sender
}

func NewCollectAndSendMetricsService(conf *flags.Config) *CollectAndSendMetricsService {
	return &CollectAndSendMetricsService{
		config: *conf,
		sender: senderPack.NewSender(conf),
	}
}

func (c *CollectAndSendMetricsService) Run() {
	var pollCount int64
	metricsRepo := repository.NewInMemoryMetricsRepository()

	// Отправляем метрики на сервер
	fmt.Println("send")
	go func() {
		for {

			c.Init(metricsRepo, pollCount)

			// metrics, err := metricsRepo.GetMetrics()
			// if err != nil {
			// 	fmt.Println("Error collecting metrics:", err)
			// 	return
			// }
			// err = c.send(metrics, pollCount)
			// if err != nil {
			// 	fmt.Println("Error sending metrics:", err)
			// }
			// atomic.AddInt64(&pollCount, 1)
			// err = c.sendPack(metrics, pollCount)
			// if err != nil {
			// 	fmt.Println("Error sending metrics:", err)
			// }

			time.Sleep(c.config.GetPollInterval())
		}

	}()

	// Собираем метрики
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

func (c *CollectAndSendMetricsService) makePack(metrics map[string]interface{}, pollCount int64) []interface{} {
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
func (c *CollectAndSendMetricsService) send(metrics map[string]interface{}, pollCount int64) error {
	var forSend = c.makePack(metrics, pollCount)
	for _, metric := range forSend {
		err := c.sender.SendMetricJSON(metric, senderPack.UpdatePathOne)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CollectAndSendMetricsService) sendPack(metrics map[string]interface{}, pollCount int64) error {
	var forSend = c.makePack(metrics, pollCount)
	err := c.sender.SendMetricJSON(forSend, senderPack.UpdatePathPack)
	if err != nil {
		return err
	}

	return nil
}
