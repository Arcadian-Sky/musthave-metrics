package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

// Config содержит конфигурационные параметры сервера
type Config struct {
	Interval        time.Duration // Интервал сохранения метрик
	FileStoragePath string        // Путь к файлу для хранения метрик
	Restore         bool          // Флаг восстановления метрик при старте сервера
}
type ConfigApp struct {
	config         Config
	MetricsStorage *storage.MemStorage
}

func NewConfig(conf *Config) *ConfigApp {
	return &ConfigApp{
		config: *conf,
	}
}

// saveMetrics сохраняет текущие значения метрик на диск с указанным интервалом
func (app *ConfigApp) SaveMetrics() {
	if app.config.FileStoragePath == "" {
		return
	}
	// Сохранение синхронно
	if app.config.Interval == 0 {
		return
	}

	ticker := time.NewTicker(app.config.Interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		app.SaveToFile()
	}
}

// saveToFile сохраняет метрики в файл
func (app *ConfigApp) SaveToFile() {

	// fmt.Printf("store.GetMetric(mType): %v\n", m.GetMetric(mType))

	allMetrics := app.MetricsStorage.GetMetrics()
	// fmt.Printf("allMetrics: %v\n", allMetrics)
	jsonData, err := json.MarshalIndent(allMetrics, "", "    ")
	if err != nil {
		return //err
	}

	// fmt.Printf("SaveToFile allMetrics: %v\n", allMetrics)
	// Записываем JSON в файл.
	err = os.WriteFile(app.config.FileStoragePath, jsonData, 0644)
	if err != nil {
		return //err
	}

	// fmt.Println("Server SaveToFile Metrics")
}

// loadMetrics загружает ранее сохраненные метрики при старте сервера
func (app *ConfigApp) LoadMetrics() error {
	if app.config.FileStoragePath == "" {
		return nil
	}
	if !app.config.Restore {
		return nil
	}

	jsonData, err := os.ReadFile(app.config.FileStoragePath)
	if err != nil {
		return nil //fmt.Errorf("ошибка чтения файла: %v", err)
	}
	var newStore map[storage.MetricType]map[string]storage.MetricValue

	err = json.Unmarshal(jsonData, &newStore)
	if err != nil {
		return fmt.Errorf("ошибка разбора JSON данных: %v", err)
	}

	// fmt.Printf("newStore: %v\n", newStore)
	// Создаем карту для хранения метрик после обработки
	processedMetrics := make(map[storage.MetricType]map[string]interface{})
	// Конвертируем MetricValue в соответствующий тип в зависимости от метрики
	for metricType, metrics := range newStore {
		processedMetrics[metricType] = make(map[string]interface{})
		for metricName, metricValue := range metrics {
			// fmt.Printf("metricValue: %v\n", metricValue)
			switch metricType {
			case storage.Gauge:
				metricValueFloat := metricValue.FloatValue
				processedMetrics[metricType][metricName] = metricValueFloat
			case storage.Counter:
				metricValueInt := metricValue.IntValue
				processedMetrics[metricType][metricName] = metricValueInt
			}
		}
	}
	// fmt.Printf("processedMetrics: %v\n", processedMetrics)
	app.MetricsStorage.SetMetrics(processedMetrics)

	// fmt.Println("Server LoadMetrics")

	return nil
}

func (app *ConfigApp) HandleRequest() {
	if app.config.Interval == 0 {
		app.SaveToFile()
	}
}
