package storage

import (
	"encoding/json"
	"fmt"
)

type Memento struct {
	metrics map[MetricType]map[string]interface{}
}

// MarshalJSON преобразует Memento в JSON.
func (m *Memento) MarshalJSON() ([]byte, error) {
	// MementoJSON является вспомогательной структурой для сериализации и десериализации Memento.
	type MementoJSON struct {
		Metrics map[MetricType]map[string]interface{} `json:"metrics"`
	}
	mementoJSON := MementoJSON{
		Metrics: m.metrics,
	}
	return json.MarshalIndent(mementoJSON, "", "    ")
}

// UnmarshalJSON преобразует JSON в Memento.
func (m *Memento) UnmarshalJSON(jsonData []byte) error {

	// MementoJSON является вспомогательной структурой для сериализации и десериализации Memento.
	type MementoJSON struct {
		Metrics map[MetricType]map[string]MetricValue `json:"metrics"`
	}

	newStore := MementoJSON{}
	err := json.Unmarshal(jsonData, &newStore)
	if err != nil {
		return fmt.Errorf("ошибка разбора JSON данных: %v", err)
	}

	// Создаем карту для хранения метрик после обработки
	// Конвертируем MetricValue в соответствующий тип в зависимости от метрики
	for metricType, metrics := range newStore.Metrics {
		for metricName, metricValue := range metrics {
			// fmt.Printf("metricValue: %v\n", metricValue)
			switch metricType {
			case Gauge:
				metricValueFloat := metricValue.FloatValue
				m.setMetric(metricType, metricName, metricValueFloat)
			case Counter:
				metricValueInt := metricValue.IntValue
				m.setMetric(metricType, metricName, metricValueInt)
			}
		}
	}

	return nil
}

// SetMetric устанавливает значение метрики в Memento
func (m *Memento) setMetric(metricType MetricType, name string, value interface{}) {
	if m.metrics == nil {
		m.metrics = make(map[MetricType]map[string]interface{})
	}
	if m.metrics[metricType] == nil {
		m.metrics[metricType] = make(map[string]interface{})
	}
	m.metrics[metricType][name] = value
}

// GetMetrics возвращает текущие метрики из хранилища == getState
func (m *Memento) GetMetrics() map[MetricType]map[string]interface{} {
	return m.metrics
}

// CreateMemento - создает Memento на основе текущего состояния storage.MemStorage
func (m *MemStorage) CreateMemento() *Memento {
	return &Memento{metrics: m.GetMetrics()}
}

// RestoreFromMemento - восстанавливает состояние storage.MemStorage из Memento
func (m *MemStorage) RestoreFromMemento(s *Memento) {
	m.SetMetrics(s.metrics)
}

// GetMetrics возвращает текущие метрики из хранилища == getState
func (m *MemStorage) GetMetrics() map[MetricType]map[string]interface{} {
	return m.metrics
}

// SetMetrics метод вызывается при инициализации для перезаписи всего хранилища == setState
func (m *MemStorage) SetMetrics(metrics map[MetricType]map[string]interface{}) {
	m.metrics = metrics
}

type MetricValue struct {
	IntValue   int64
	FloatValue float64
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler
func (mv *MetricValue) UnmarshalJSON(data []byte) error {
	var intValue int64
	if err := json.Unmarshal(data, &intValue); err == nil {
		mv.IntValue = intValue
		return nil
	}

	var floatValue float64
	if err := json.Unmarshal(data, &floatValue); err != nil {
		return err
	}
	mv.FloatValue = floatValue
	return nil
}
