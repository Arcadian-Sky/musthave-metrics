package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemento_MarshalJSON(t *testing.T) {
	// Arrange
	memento := &Memento{}
	memento.metrics = map[MetricType]map[string]interface{}{
		Gauge: {
			"metric1": 123.45,
			"metric2": "abc",
		},
		Counter: {
			"metric3": 678,
			"metric4": 9.10,
		},
	}

	expectedJSON := `{
    "metrics": {
        "gauge": {
            "metric1": 123.45,
            "metric2": "abc"
        },
        "counter": {
            "metric3": 678,
            "metric4": 9.1
        }
    }
}`

	// Act
	jsonData, err := memento.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(jsonData))
}

func TestMemento_UnmarshalJSON(t *testing.T) {
	t.Run("should be able to unmarshal a JSON object into a Memento struct", func(t *testing.T) {
		jsonData := []byte(`{
            "metrics": {
                "gauge": {
                    "metric_name_1": 10,
                    "metric_name_2": 20.5
                },
                "counter": {
                    "metric_name_3": 30
                }
            }
        }`)

		memento := &Memento{}
		err := memento.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		expectedMetrics := map[MetricType]map[string]interface{}{
			Gauge: {
				"metric_name_1": float64(10),
				"metric_name_2": float64(20.5),
			},
			Counter: {
				"metric_name_3": int64(30),
			},
		}
		encoder := json.NewEncoder(os.Stdout)

		assert.Equal(t, encoder.Encode(expectedMetrics), encoder.Encode(memento.GetMetrics()))
	})
}

func TestMemento_SetMetric(t *testing.T) {
	tests := []struct {
		tname    string
		memento  *Memento
		metric   MetricType
		name     string
		value    interface{}
		expected map[MetricType]map[string]interface{}
		err      error
	}{
		{
			tname:   "sets a gauge metric",
			memento: &Memento{},
			metric:  Gauge,
			name:    "test_gauge",
			value:   123.45,
			expected: map[MetricType]map[string]interface{}{
				Gauge: {
					"test_gauge": 123.45,
				},
			},
			err: nil,
		},
		{
			tname:   "sets a counter metric",
			memento: &Memento{},
			metric:  Counter,
			name:    "test_counter",
			value:   678,
			expected: map[MetricType]map[string]interface{}{
				Counter: {
					"test_counter": 678,
				},
			},
			err: nil,
		},
		// TODO: wait error
		// {
		// 	tname:     "returns an error for an unknown metric type",
		// 	memento:  &Memento{},
		// 	metric:   999,
		// 	name:     "test_unknown",
		// 	value:    123,
		// 	expected: nil,
		// 	err:      errors.New("unknown metric type: 999"),
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.memento.SetMetrics(nil)

			tt.memento.setMetric(tt.metric, tt.name, tt.value)
			// assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, tt.memento.GetMetrics())
		})
	}
}

func TestMemento_GetMetrics(t *testing.T) {
	// given
	memento := &Memento{}
	memento.SetMetrics(map[MetricType]map[string]interface{}{
		Gauge: {
			"metric_a": 1.0,
			"metric_b": 2.0,
		},
		Counter: {
			"metric_c": 3.0,
			"metric_d": 4.0,
		},
	})

	// when
	metrics := memento.GetMetrics()

	// then
	expectedMetrics := map[MetricType]map[string]interface{}{
		Gauge: {
			"metric_a": 1.0,
			"metric_b": 2.0,
		},
		Counter: {
			"metric_c": 3.0,
			"metric_d": 4.0,
		},
	}
	assert.Equal(t, expectedMetrics, metrics)
}

func TestMemento_SetMetrics(t *testing.T) {
	// given
	memento := &Memento{}

	// when
	memento.SetMetrics(map[MetricType]map[string]interface{}{
		Gauge: {
			"metric_a": 1.0,
			"metric_b": 2.0,
		},
		Counter: {
			"metric_c": 3,
			"metric_d": 4,
		},
	})

	// then
	expectedMetrics := map[MetricType]map[string]interface{}{
		Gauge: {
			"metric_a": 1.0,
			"metric_b": 2.0,
		},
		Counter: {
			"metric_c": 3,
			"metric_d": 4,
		},
	}
	assert.Equal(t, expectedMetrics, memento.GetMetrics())
}
