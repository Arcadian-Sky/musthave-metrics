package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config содержит конфигурационные параметры сервера
type Config struct {
	Interval        time.Duration // Интервал сохранения метрик
	FileStoragePath string        // Путь к файлу для хранения метрик
	Restore         bool          // Флаг восстановления метрик при старте сервера
}

// saveMetrics сохраняет текущие значения метрик на диск с указанным интервалом
func (m *MemStorage) SaveMetrics(config Config) {
	if config.Interval == 0 {
		// Сохранение синхронно
		m.SaveToFile(config.FileStoragePath, m)
		return
	}

	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.SaveToFile(config.FileStoragePath, m)
		}
	}
}

// saveToFile сохраняет метрики в файл
func (m *MemStorage) SaveToFile(filePath string) {
	// mType, _ := GetMetricTypeByCode("gauge")
	// fmt.Printf("store.GetMetric(mType): %v\n", m.GetMetric(mType))

	allMetrics := m.GetMetrics()
	jsonData, err := json.MarshalIndent(allMetrics, "", "    ")
	if err != nil {
		return //err
	}

	// Записываем JSON в файл.
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return //err
	}
}

// loadMetrics загружает ранее сохраненные метрики при старте сервера
func (m *MemStorage) LoadMetrics(config Config) error {
	jsonData, err := os.ReadFile(config.FileStoragePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}
	var newStore map[MetricType]map[string]interface{}
	err = json.Unmarshal(jsonData, &newStore)
	if err != nil {
		return fmt.Errorf("ошибка разбора JSON данных: %v", err)
	}
	m.SetMetrics(newStore)
	return nil
}
