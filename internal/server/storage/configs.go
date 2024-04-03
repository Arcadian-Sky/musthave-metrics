package storage

import (
	"encoding/json"
)

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
