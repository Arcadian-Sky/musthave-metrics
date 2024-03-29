package models

// В теле ответа отправляйте JSON той же структуры с актуальным (изменённым) значением Value(/Delta).
// В теле запроса должен быть описанный выше JSON с заполненными полями ID и MType
// Metrics структура для передачи метрик.
// @Summary Метрики
// @Description Структура для передачи метрик
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
