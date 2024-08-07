package controller

import (
	"math"
	"sync"
	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/repository"
)

// MockMetricsRepository представляет макет для тестирования.
type MockMetricsRepository struct {
	metrics map[string]interface{}
}

// GetMetrics возвращает метрики и ошибку (если есть).
func (m *MockMetricsRepository) GetMetrics() (map[string]interface{}, error) {
	return m.metrics, nil
}

func TestCollectAndSendMetricsService_GetSystemInfo(t *testing.T) {
	// Создаем экземпляр репозитория метрик
	mockMetricsRepo := &MockMetricsRepository{
		metrics: map[string]interface{}{
			"ExistingMetric": 123,
		},
	}

	// Создаем канал для передачи метрик
	metricsInfoChan := make(chan map[string]interface{})

	var wg sync.WaitGroup
	wg.Add(1)

	// Запускаем функцию getSystemInfo в отдельной горутине
	go func() {
		// Запускаем функцию getSystemInfo в отдельной горутине
		getSystemInfo(metricsInfoChan, mockMetricsRepo)
		wg.Done()
	}()

	// Получаем результат из канала
	receivedMetrics := <-metricsInfoChan

	// Ожидаем завершения выполнения горутины
	wg.Wait()

	// Проверяем, что полученные метрики соответствуют ожидаемым
	assert.Equal(t, 123, receivedMetrics["ExistingMetric"])
	// Проверяем другие ожидаемые метрики

	memoryInfo, _ := mem.VirtualMemory()
	cpuCount, _ := cpu.Counts(false)

	numberRTotal := int(math.Round(receivedMetrics["TotalMemory"].(float64)) / 1000000000)
	numberTotal := int(math.Round(float64(memoryInfo.Total)) / 1000000000)

	numberRFree := int(math.Round(receivedMetrics["FreeMemory"].(float64)) / 1000000000)
	numberFree := int(math.Round(float64(memoryInfo.Free)) / 1000000000)

	numberRcpuCount := int(math.Round(receivedMetrics["CPUutilization1"].(float64)) / 1000000000)
	numbercpuCount := int(math.Round(float64(cpuCount)) / 1000000000)

	// Проверяем, что TotalMemory, FreeMemory и CPUutilization1 установлены относительно правильно
	assert.Equal(t, numberTotal, numberRTotal)
	assert.Equal(t, numberFree, numberRFree)
	assert.Equal(t, numbercpuCount, numberRcpuCount)

}

func TestCollectAndSendMetricsService_Init(t *testing.T) {
	// Создаем экземпляр репозитория метрик
	metricsRepo := repository.NewInMemoryMetricsRepository()

	// Создаем экземпляр сервиса, который будем тестировать
	service := NewCollectAndSendMetricsService(flags.SetDefault())

	var pollCount int64 = 5
	// Инициализируем сервис
	service.Init(metricsRepo, &pollCount)

	// Ожидаем, что количество горутин воркеров будет равно количеству воркеров
	// expectedWorkers := 5
	// time.Sleep(time.Millisecond)
	// actualWorkers := runtime.NumGoroutine() // Получаем количество горутин

	// if actualWorkers != expectedWorkers+1 {
	// 	// +1, потому что есть главная горутина исполнения теста
	// 	t.Errorf("Expected %d workers goroutines, but got %d", expectedWorkers, actualWorkers-1)
	// }
}

func TestCollectAndSendMetricsService_Push(t *testing.T) {
	// type fields struct {
	// 	config flags.Config
	// 	sender *senderPack.Sender
	// }
	// type args struct {
	// 	metrics   map[string]interface{}
	// 	pollCount int64
	// }
	// tests := []struct {
	// 	name   string
	// 	fields fields
	// 	args   args
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		c := &CollectAndSendMetricsService{
	// 			config: tt.fields.config,
	// 			sender: tt.fields.sender,
	// 		}
	// 		c.Push(tt.args.metrics, tt.args.pollCount)
	// 	})
	// }

	metricsRepo := repository.NewInMemoryMetricsRepository()

	// Создаем экземпляр сервиса, который будем тестировать
	service := NewCollectAndSendMetricsService(flags.SetDefault())

	var pollCount int64 = 5

	metrics, _ := metricsRepo.GetMetrics()

	// Инициализируем сервис
	service.Push(metrics, &pollCount)
}
