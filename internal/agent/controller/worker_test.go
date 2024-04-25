package controller

import (
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/repository"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"
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

	// Запускаем функцию getSystemInfo в отдельной горутине
	go getSystemInfo(metricsInfoChan, mockMetricsRepo)

	// Получаем результат из канала
	receivedMetrics := <-metricsInfoChan

	// Проверяем, что полученные метрики соответствуют ожидаемым
	assert.Equal(t, 123, receivedMetrics["ExistingMetric"])
	// Проверяем другие ожидаемые метрики

	memoryInfo, _ := mem.VirtualMemory()
	cpuCount, _ := cpu.Counts(false)

	// Проверяем, что TotalMemory, FreeMemory и CPUutilization1 установлены правильно
	assert.Equal(t, float64(memoryInfo.Total), receivedMetrics["TotalMemory"])
	assert.Equal(t, float64(memoryInfo.Free), receivedMetrics["FreeMemory"])
	assert.Equal(t, float64(cpuCount), receivedMetrics["CPUutilization1"])

}

func TestCollectAndSendMetricsService_Init(t *testing.T) {
	// Создаем экземпляр репозитория метрик
	metricsRepo := repository.NewInMemoryMetricsRepository()

	// Создаем экземпляр сервиса, который будем тестировать
	service := NewCollectAndSendMetricsService(flags.SetDefault())

	// Инициализируем сервис
	service.Init(metricsRepo, 5)

	// Ожидаем, что количество горутин воркеров будет равно количеству воркеров
	// expectedWorkers := 5
	// time.Sleep(time.Millisecond)
	// actualWorkers := runtime.NumGoroutine() // Получаем количество горутин

	// if actualWorkers != expectedWorkers+1 {
	// 	// +1, потому что есть главная горутина исполнения теста
	// 	t.Errorf("Expected %d workers goroutines, but got %d", expectedWorkers, actualWorkers-1)
	// }
}
