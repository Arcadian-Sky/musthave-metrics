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
func (m *MemStorage) SaveMetrics() {
	if m.conf.FileStoragePath == "" {
		return
	}
	// Сохранение синхронно
	if m.conf.Interval == 0 {
		return
	}

	ticker := time.NewTicker(m.conf.Interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		m.SaveToFile()
	}
}

// saveToFile сохраняет метрики в файл
func (m *MemStorage) SaveToFile() {

	// fmt.Printf("store.GetMetric(mType): %v\n", m.GetMetric(mType))

	allMetrics := m.GetMetrics()
	jsonData, err := json.MarshalIndent(allMetrics, "", "    ")
	if err != nil {
		return //err
	}

	// fmt.Printf("SaveToFile allMetrics: %v\n", allMetrics)
	// Записываем JSON в файл.
	err = os.WriteFile(m.conf.FileStoragePath, jsonData, 0644)
	if err != nil {
		return //err
	}

	fmt.Println("Server SaveToFile Metrics")
}

// loadMetrics загружает ранее сохраненные метрики при старте сервера
func (m *MemStorage) LoadMetrics() error {
	if m.conf.FileStoragePath == "" {
		return nil
	}
	if !m.conf.Restore {
		return nil
	}

	jsonData, err := os.ReadFile(m.conf.FileStoragePath)
	if err != nil {
		return nil //fmt.Errorf("ошибка чтения файла: %v", err)
	}
	var newStore map[MetricType]map[string]interface{}
	err = json.Unmarshal(jsonData, &newStore)
	if err != nil {
		return fmt.Errorf("ошибка разбора JSON данных: %v", err)
	}
	m.SetMetrics(newStore)
	// fmt.Printf("newStore: %v\n", newStore)
	// mType, _ := GetMetricTypeByCode("counter")
	// metrics := m.GetMetric(mType)
	// metrics := m.GetMetrics()
	// for name, value := range metrics {
	// 	fmt.Printf("LoadMetrics name: %v value: %v\n", name, value)
	// }

	fmt.Println("Server LoadMetrics")

	return nil
}
