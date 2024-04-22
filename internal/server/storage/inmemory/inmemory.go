package inmemory

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/utils"
)

// MemStorage представляет хранилище метрик
type MemStorage struct {
	metrics map[storage.MetricType]map[string]interface{}
}

// NewMemStorage создает новый экземпляр MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[storage.MetricType]map[string]interface{}),
	}
}

func (m *MemStorage) GetJSONMetric(ctx context.Context, metric *models.Metrics) error {
	metricType, err := utils.GetMetricTypeByCode(metric.MType)

	if err != nil {
		return err
	}

	if _, ok := m.metrics[metricType]; !ok {
		m.metrics[metricType] = make(map[string]interface{})
	}
	realVal := m.metrics[metricType][metric.ID]
	switch metricType {
	case storage.Gauge:
		if f, ok := realVal.(float64); ok {
			metric.Value = &f
		} else {
			zeroValue := float64(0)
			metric.Value = &zeroValue
		}
	case storage.Counter:
		if i, ok := realVal.(int64); ok {
			metric.Delta = &i
		} else {
			zeroValue := int64(0)
			metric.Delta = &zeroValue
		}
		fmt.Printf("Get saved metric.Delta: %v\n", realVal)
		fmt.Printf("Get returned metric.Delta: %v\n", *metric.Delta)
	default:
		return fmt.Errorf("invalid metric type")
	}

	return nil
}

func (m *MemStorage) GetJSONMetrics(ctx context.Context, metric *[]models.Metrics) error {
	// metricType, err := utils.GetMetricTypeByCode(metric.MType)

	// if err != nil {
	// 	return err
	// }

	// if _, ok := m.metrics[metricType]; !ok {
	// 	m.metrics[metricType] = make(map[string]interface{})
	// }
	// realVal := m.metrics[metricType][metric.ID]
	// switch metricType {
	// case storage.Gauge:
	// 	if f, ok := realVal.(float64); ok {
	// 		metric.Value = &f
	// 	} else {
	// 		zeroValue := float64(0)
	// 		metric.Value = &zeroValue
	// 	}
	// case storage.Counter:
	// 	if i, ok := realVal.(int64); ok {
	// 		metric.Delta = &i
	// 	} else {
	// 		zeroValue := int64(0)
	// 		metric.Delta = &zeroValue
	// 	}
	// 	// fmt.Printf("Get saved metric.Delta: %v\n", realVal)
	// 	// fmt.Printf("Get returned metric.Delta: %v\n", *metric.Delta)
	// default:
	// 	return fmt.Errorf("invalid metric type")
	// }

	return nil
}

func (m *MemStorage) UpdateJSONMetric(ctx context.Context, metric *models.Metrics) error {
	metricType, err := utils.GetMetricTypeByCode(metric.MType)

	if err != nil {
		return err
	}
	if m.metrics == nil {
		m.metrics = make(map[storage.MetricType]map[string]interface{})
	}
	if _, ok := m.metrics[metricType]; !ok {
		m.metrics[metricType] = make(map[string]interface{})
	}
	switch metricType {
	case storage.Gauge:
		if metric.Value == nil {
			zeroValue := float64(0)
			metric.Value = &zeroValue
		}
		m.metrics[metricType][metric.ID] = *metric.Value
	case storage.Counter:
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
	default:
		return fmt.Errorf("invalid metric type")
	}

	return nil
}

// UpdateMetric обновляет значение метрики в хранилище
func (m *MemStorage) UpdateMetric(ctx context.Context, mtype string, name string, value string) error {
	// var metricType MetricType
	metricType, err := utils.GetMetricTypeByCode(mtype)
	if err != nil {
		return err
	}

	if _, ok := m.metrics[metricType]; !ok {
		m.metrics[metricType] = make(map[string]interface{})
	}

	switch metricType {
	case storage.Gauge:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			m.metrics[metricType][name] = floatValue
		} else {
			return fmt.Errorf("invalid metric value: %v", err)
		}
	case storage.Counter:
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

func (m *MemStorage) UpdateJSONMetrics(ctx context.Context, metrics *[]models.Metrics) error {
	// var metricType MetricType

	return nil
}

// GetMetric возвращает текущие метрики из хранилища для типа
func (m *MemStorage) GetMetric(ctx context.Context, mtype storage.MetricType) map[string]interface{} {
	return m.metrics[mtype]
}

// GetMetrics возвращает текущие метрики из хранилища == getState
func (m *MemStorage) GetMetrics(ctx context.Context) map[storage.MetricType]map[string]interface{} {
	return m.metrics
}

// SetMetrics метод вызывается при инициализации для перезаписи всего хранилища == setState
func (m *MemStorage) SetMetrics(ctx context.Context, metrics map[storage.MetricType]map[string]interface{}) {
	m.metrics = metrics
}

func (m *MemStorage) Ping() error {
	return fmt.Errorf("формат хранения не поддерживает бд")
}

// CreateMemento - создает Memento на основе текущего состояния storage.MemStorage
func (m *MemStorage) CreateMemento() *storage.Memento {
	s := &storage.Memento{}
	ctx := context.Background()
	s.SetMetrics(m.GetMetrics(ctx))
	return s
}

// RestoreFromMemento - восстанавливает состояние storage.MemStorage из Memento
func (m *MemStorage) RestoreFromMemento(s *storage.Memento) {
	ctx := context.Background()
	m.SetMetrics(ctx, s.GetMetrics())
}
