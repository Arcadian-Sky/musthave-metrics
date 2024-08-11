package grpc

import (
	"testing"

	pb "github.com/Arcadian-Sky/musthave-metrics/gen/proto/api/metrics/v1"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestConvertMetrics(t *testing.T) {
	metrics := map[storage.MetricType]map[string]interface{}{
		storage.Gauge: {
			"metric1": 1.23,
		},
		storage.Counter: {
			"metric2": int64(10),
		},
	}

	expected := []*pb.Metric{
		{
			Id:   "metric1",
			Type: pb.Type_TYPE_GAUGE,
			Mvalue: &pb.Metric_Value{
				Value: 1.23,
			},
		},
		{
			Id:   "metric2",
			Type: pb.Type_TYPE_COUNTER,
			Mvalue: &pb.Metric_Delta{
				Delta: 10,
			},
		},
	}

	result := convertMetrics(metrics)
	assert.ElementsMatch(t, expected, result, "The converted metrics do not match the expected values")
}

func TestConvertModelMetricsToRPCMetrics(t *testing.T) {
	tests := []struct {
		name     string
		input    *models.Metrics
		expected *pb.Metric
	}{
		{
			name: "Gauge Metric with Value",
			input: &models.Metrics{
				ID:    "metric1",
				MType: "gauge",
				Value: float64Ptr(1.23),
			},
			expected: &pb.Metric{
				Id:   "metric1",
				Type: pb.Type_TYPE_GAUGE,
				Mvalue: &pb.Metric_Value{
					Value: 1.23,
				},
			},
		},
		{
			name: "Counter Metric with Delta",
			input: &models.Metrics{
				ID:    "metric2",
				MType: "counter",
				Delta: int64Ptr(10),
			},
			expected: &pb.Metric{
				Id:   "metric2",
				Type: pb.Type_TYPE_COUNTER,
				Mvalue: &pb.Metric_Delta{
					Delta: 10,
				},
			},
		},
		{
			name: "Metric with Nil Value and Delta",
			input: &models.Metrics{
				ID:    "metric3",
				MType: "gauge",
				Value: nil,
				Delta: int64Ptr(5),
			},
			expected: &pb.Metric{
				Id:     "metric3",
				Type:   pb.Type_TYPE_GAUGE,
				Mvalue: &pb.Metric_Value{},
			},
		},
		{
			name: "Metric with Unknown Type",
			input: &models.Metrics{
				ID:    "metric4",
				MType: "unknown",
				Value: float64Ptr(3.14),
			},
			expected: &pb.Metric{
				Id:   "metric4",
				Type: pb.Type_TYPE_UNSPECIFIED,
			},
		},
		{
			name:     "Metric with Empty ID",
			input:    &models.Metrics{ID: "", MType: "gauge", Value: float64Ptr(2.71)},
			expected: &pb.Metric{Id: "", Type: pb.Type_TYPE_GAUGE, Mvalue: &pb.Metric_Value{Value: 2.71}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := convertModelMetricsToRPCMetrics(test.input)
			assert.Equal(t, test.expected, result, "The converted model metrics do not match the expected values")
		})
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}

func TestConvertStringToMetricType(t *testing.T) {
	tests := []struct {
		input    string
		expected pb.Type
	}{
		{string(storage.Gauge), pb.Type_TYPE_GAUGE},
		{string(storage.Counter), pb.Type_TYPE_COUNTER},
		{"unknown", pb.Type_TYPE_UNSPECIFIED},
	}

	for _, test := range tests {
		result := convertStringToMetricType(test.input)
		assert.Equal(t, test.expected, result, "The converted metric type does not match the expected value")
	}
}

func TestConvertMetricTypeToString(t *testing.T) {
	tests := []struct {
		input    pb.Type
		expected string
	}{
		{pb.Type_TYPE_GAUGE, string(storage.Gauge)},
		{pb.Type_TYPE_COUNTER, string(storage.Counter)},
		{pb.Type_TYPE_UNSPECIFIED, ""},
	}

	for _, test := range tests {
		result := convertMetricTypeToString(test.input)
		assert.Equal(t, test.expected, result, "The converted metric type string does not match the expected value")
	}
}
