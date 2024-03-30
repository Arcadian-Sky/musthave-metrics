package storage

import (
	"fmt"
	"strconv"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
)

// MetricType определяет тип метрики (gauge или counter)
type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

// MetricsStorage определяет интерфейс для взаимодействия с хранилищем метрик
type MetricsStorage interface {
	GetJSONMetric(metric *models.Metrics) error
	UpdateJSONMetric(metric *models.Metrics) error
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

func (m *MemStorage) GetJSONMetric(metric *models.Metrics) error {
	metricType, err := GetMetricTypeByCode(metric.MType)

	if err != nil {
		return err
	}

	if _, ok := m.metrics[metricType]; !ok {
		m.metrics[metricType] = make(map[string]interface{})
	}
	realVal := m.metrics[metricType][metric.ID]
	// fmt.Printf("realVal: %v\n", realVal)
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

	if _, ok := m.metrics[metricType]; !ok {
		m.metrics[metricType] = make(map[string]interface{})
	}

	switch metricType {
	case Gauge:
		if metric.Value == nil {
			return fmt.Errorf("value must be defined")
		}
		m.metrics[metricType][metric.ID] = metric.Value
	case Counter:
		if metric.Delta == nil {
			return fmt.Errorf("delta must be defined")
		}
		currentCounter, ok := m.metrics[metricType][metric.ID].(int64)
		if ok {
			m.metrics[metricType][metric.ID] = currentCounter + *metric.Delta
		} else {
			m.metrics[metricType][metric.ID] = *metric.Delta
		}
	default:
		return fmt.Errorf("invalid metric type")
	}

	fmt.Println(metricType)

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

	return nil
}

// GetMetrics возвращает текущие метрики из хранилища
func (m *MemStorage) GetMetrics() map[MetricType]map[string]interface{} {
	return m.metrics
}

// GetMetric возвращает текущие метрики из хранилища для типа
func (m *MemStorage) GetMetric(mtype MetricType) map[string]interface{} {
	return m.metrics[mtype]
}
