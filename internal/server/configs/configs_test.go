package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

// TestSaveToFile проверяет метод SaveToFile, убеждаясь, что метрики сохраняются в файле JSON.
func TestSaveToFile(t *testing.T) {
	// Создаем временный файл для сохранения метрик
	tmpfile, err := ioutil.TempFile("", "test_metrics.json")
	if err != nil {
		t.Fatalf("ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Удаляем временный файл после завершения теста

	// Создаем экземпляр ConfigApp с настройками и хранилищем метрик
	conf := Config{
		Interval:        1 * time.Second, // Устанавливаем небольшой интервал для теста
		FileStoragePath: tmpfile.Name(),  // Используем временный файл для сохранения метрик
		Restore:         false,           // Не восстанавливаем метрики при старте
	}
	metricsStorage := storage.NewMemStorage()
	app := NewConfig(&conf)
	app.MetricsStorage = metricsStorage

	// Добавляем какие-то метрики в хранилище для сохранения
	metrics := map[storage.MetricType]map[string]interface{}{
		storage.Gauge: {
			"metric1": float64(10.5),
			"metric2": float64(20.8),
		},
		storage.Counter: {
			"metric3": int64(100),
			"metric4": int64(200),
		},
	}
	metricsStorage.SetMetrics(metrics)

	// Вызываем метод SaveToFile для сохранения метрик
	app.SaveToFile()

	// Читаем содержимое временного файла
	fileContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ошибка чтения временного файла: %v", err)
	}

	// Проверяем, что содержимое файла соответствует ожидаемым метрикам
	var savedMetrics map[storage.MetricType]map[string]storage.MetricValue
	if err := json.Unmarshal(fileContent, &savedMetrics); err != nil {
		t.Fatalf("ошибка разбора JSON данных из файла: %v", err)
	}
	fmt.Printf("savedMetrics: %v\n", savedMetrics)
	fmt.Printf("metrics: %v\n", metrics)
	// Сравниваем ожидаемые метрики с сохраненными в файле
	if !compareMetricMaps(metrics, savedMetrics) {
		t.Error("сохраненные метрики не совпадают с ожидаемыми")
	}
}

// compareMetricMaps сравнивает два набора метрик на равенство.
func compareMetricMaps(expected map[storage.MetricType]map[string]interface{}, actual map[storage.MetricType]map[string]storage.MetricValue) bool {
	if len(expected) != len(actual) {
		return false
	}
	for metricType, expectedMetrics := range expected {
		actualMetrics, ok := actual[metricType]
		if !ok || len(expectedMetrics) != len(actualMetrics) {
			fmt.Printf("\"1err\": %v\n", len(expectedMetrics) != len(actualMetrics))
			return false
		}
		for metricName, expectedValue := range expectedMetrics {
			actualValue, ok := actualMetrics[metricName]
			if !ok {
				return false
			}
			if reflect.TypeOf(expectedValue) == reflect.TypeOf(float64(0)) {
				actualFloatValue := actualValue.FloatValue
				if expectedValue != actualFloatValue {
					return false
				}
			} else if reflect.TypeOf(expectedValue) == reflect.TypeOf(int64(0)) {
				actualFloatValue := actualValue.IntValue
				if expectedValue != actualFloatValue {
					return false
				}
			}
		}
	}
	return true
}

func TestSaveMetrics(t *testing.T) {
	// Создание временного файла для тестирования
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Создание экземпляра ConfigApp с настройками для тестирования
	config := Config{
		Interval:        100 * time.Millisecond, // Небольшой интервал для быстрого завершения теста
		FileStoragePath: tmpfile.Name(),
		Restore:         true,
	}
	app := NewConfig(&config)
	app.MetricsStorage = storage.NewMemStorage()
	// Вызов SaveMetrics
	go app.SaveMetrics()

	// Ждем некоторое время, чтобы прошло несколько циклов сохранения
	time.Sleep(500 * time.Millisecond)

	// Проверяем, что файл был создан и не пустой
	fileInfo, err := os.Stat(tmpfile.Name())
	if err != nil {
		t.Fatalf("error getting file info: %v", err)
	}
	assert.False(t, fileInfo.Size() == 0, "expected file not to be empty")

	// Отключаем цикл сохранения для завершения теста
	tmpfile.Close()
}

func TestLoadMetrics(t *testing.T) {
	// Создаем временный файл с данными метрик для теста
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	// defer os.Remove(tmpfile.Name())

	// Создаем конфигурацию для теста
	conf := Config{
		Interval:        time.Second * 10,
		FileStoragePath: tmpfile.Name(),
		Restore:         true,
	}
	app := NewConfig(&conf)
	metricStorege := storage.NewMemStorage()
	app.MetricsStorage = metricStorege

	// Создаем тестовые данные для метрик
	testMetrics := map[storage.MetricType]map[string]interface{}{
		storage.Gauge: {
			"metric1": 1.23,
			"metric2": 4.56,
		},
		storage.Counter: {
			"metric3": 100,
			"metric4": 200,
		},
	}

	// Сохраняем тестовые данные в файл
	jsonData, err := json.Marshal(testMetrics)
	assert.NoError(t, err)
	err = os.WriteFile(tmpfile.Name(), jsonData, 0644)
	assert.NoError(t, err)

	// Загружаем метрики из файла
	err = app.LoadMetrics()
	assert.NoError(t, err)

	// Проверяем, что метрики были успешно загружены и сохранены в MemStorage
	expectedJSON, err := json.Marshal(testMetrics)
	assert.NoError(t, err)

	actualJSON, err := json.Marshal(metricStorege.GetMetrics())
	assert.NoError(t, err)

	assert.Equal(t, expectedJSON, actualJSON)
}
