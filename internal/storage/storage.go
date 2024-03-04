package storage

import (
	"fmt"
	"strconv"
)

var Storage *MemStorage

// MetricType определяет тип метрики (gauge или counter)
type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

// MetricsStorage определяет интерфейс для взаимодействия с хранилищем метрик
type MetricsStorage interface {
	UpdateMetric(mtype string, name string, value string) error
	GetMetrics() map[MetricType]map[string]interface{}
	GetMetric(MetricType) map[string]interface{}
}

// MemStorage представляет хранилище метрик
type MemStorage struct {
	metrics map[MetricType]map[string]interface{}
}

// NewMemStorage создает новый экземпляр MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[MetricType]map[string]interface{}),
	}
}

// UpdateMetric обновляет значение метрики в хранилище
func (m *MemStorage) UpdateMetric(mtype string, name string, value string) error {

	var metricType MetricType
	switch mtype {
	case "gauge":
		metricType = Gauge
	case "counter":
		metricType = Counter
	default:
		return fmt.Errorf("invalid metric type")
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

	return nil
}

// GetMetrics возвращает текущие метрики из хранилища
func (m *MemStorage) GetMetrics() map[MetricType]map[string]interface{} {
	return m.metrics
}

func (m *MemStorage) GetMetric(mtype MetricType) map[string]interface{} {
	return m.metrics[mtype]
}
