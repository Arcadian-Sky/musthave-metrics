package storage

import (
	"fmt"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
)

// MetricType определяет тип метрики (gauge или counter)
type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

// MetricsStorage определяет интерфейс для взаимодействия с хранилищем метрик
type MetricsStorage interface {
	GetMetric(mtype MetricType) map[string]interface{}
	UpdateMetric(mtype string, name string, value string) error

	GetJSONMetric(metric *models.Metrics) error
	UpdateJSONMetric(metric *models.Metrics) error

	// GetJSONMetrics(metrics *[]models.Metrics)
	UpdateJSONMetrics(metrics *[]models.Metrics) error

	GetMetrics() map[MetricType]map[string]interface{}
	SetMetrics(metrics map[MetricType]map[string]interface{})

	Ping() error
}

func GetMetricTypeByCode(mtype string) (MetricType, error) {
	var metricType MetricType
	switch mtype {
	case "gauge":
		metricType = Gauge
	case "counter":
		metricType = Counter
	default:
		return metricType, fmt.Errorf("invalid metric type")
	}

	return metricType, nil
}

/*
// MemStorage представляет хранилище метрик
type MemStorage struct {
	metrics map[MetricType]map[string]interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[MetricType]map[string]interface{}),
	}
}

func (m *MemStorage) GetJSONMetric(metric *models.Metrics) error {
	metricType, err := GetMetricTypeByCode(metric.MType)

	if err != nil {
		return err
	}

	if _, ok := m.metrics[metricType]; !ok {
		m.metrics[metricType] = make(map[string]interface{})
	}
	realVal := m.metrics[metricType][metric.ID]
	switch metricType {
	case Gauge:
		if f, ok := realVal.(float64); ok {
			metric.Value = &f
		} else {
			zeroValue := float64(0)
			metric.Value = &zeroValue
		}
	case Counter:
		if i, ok := realVal.(int64); ok {
			metric.Delta = &i
		} else {
			zeroValue := int64(0)
			metric.Delta = &zeroValue
		}
		// fmt.Printf("Get saved metric.Delta: %v\n", realVal)
		// fmt.Printf("Get returned metric.Delta: %v\n", *metric.Delta)
	default:
		return fmt.Errorf("invalid metric type")
	}

	return nil
}

func (m *MemStorage) UpdateJSONMetric(metric *models.Metrics) error {
	metricType, err := GetMetricTypeByCode(metric.MType)

	if err != nil {
		return err
	}
	if m.metrics == nil {
		m.metrics = make(map[MetricType]map[string]interface{})
	}
	if _, ok := m.metrics[metricType]; !ok {
		m.metrics[metricType] = make(map[string]interface{})
	}
	switch metricType {
	case Gauge:
		if metric.Value == nil {
			zeroValue := float64(0)
			metric.Value = &zeroValue
		}
		m.metrics[metricType][metric.ID] = *metric.Value
	case Counter:
		// fmt.Printf("Update val metric.Delta: %v\n", *metric.Delta)
		if metric.Delta == nil {
			zeroValue := int64(0)
			metric.Delta = &zeroValue
		}
		currentCounter, ok := m.metrics[metricType][metric.ID].(int64)
		if ok {
			m.metrics[metricType][metric.ID] = currentCounter + *metric.Delta
		} else {
			m.metrics[metricType][metric.ID] = *metric.Delta
		}
		// fmt.Printf("Update saved metric.Delta: %v\n", m.metrics[metricType][metric.ID])
	default:
		return fmt.Errorf("invalid metric type")
	}

	// if m.conf.Interval == 0 {
	// 	// m.SaveToFile()
	// }

	return nil
}

// UpdateMetric обновляет значение метрики в хранилище
func (m *MemStorage) UpdateMetric(mtype string, name string, value string) error {

	// var metricType MetricType
	metricType, err := GetMetricTypeByCode(mtype)

	if err != nil {
		return err
	}

	if _, ok := m.metrics[metricType]; !ok {
		m.metrics[metricType] = make(map[string]interface{})
	}

	switch metricType {
	case Gauge:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			m.metrics[metricType][name] = floatValue
		} else {
			return fmt.Errorf("invalid metric value: %v", err)
		}
	case Counter:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			currentCounter, ok := m.metrics[metricType][name].(int64)
			if ok {
				m.metrics[metricType][name] = currentCounter + intValue
			} else {
				m.metrics[metricType][name] = intValue
			}
		} else {
			return fmt.Errorf("invalid metric value: %v", err)
		}
	default:
		return fmt.Errorf("invalid metric type")
	}

	// if m.conf.Interval == 0 {
	// 	// m.SaveToFile()
	// }

	return nil
}

// GetMetric возвращает текущие метрики из хранилища для типа
func (m *MemStorage) GetMetric(mtype MetricType) map[string]interface{} {
	return m.metrics[mtype]
}
*/
