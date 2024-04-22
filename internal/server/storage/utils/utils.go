package utils

import (
	"fmt"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

func GetMetricTypeByCode(mtype string) (storage.MetricType, error) {
	var metricType storage.MetricType
	switch mtype {
	case "gauge":
		metricType = storage.Gauge
	case "counter":
		metricType = storage.Counter
	default:
		return metricType, fmt.Errorf("invalid metric type")
	}

	return metricType, nil
}
