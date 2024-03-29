package storage

import (
	"reflect"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
)

func TestMemStorage_UpdateMetric(t *testing.T) {
	type fields struct {
		metrics map[MetricType]map[string]interface{}
	}
	type args struct {
		mtype string
		name  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UpdateGaugeMetric",
			fields: fields{
				metrics: map[MetricType]map[string]interface{}{},
			},
			args: args{
				mtype: "gauge",
				name:  "testGauge",
				value: "123",
			},
			wantErr: false,
		},
		{
			name: "UpdateCounterMetric",
			fields: fields{
				metrics: map[MetricType]map[string]interface{}{},
			},
			args: args{
				mtype: "counter",
				name:  "testCounter",
				value: "456",
			},
			wantErr: false,
		},
		{
			name: "InvalidMetricType",
			fields: fields{
				metrics: map[MetricType]map[string]interface{}{},
			},
			args: args{
				mtype: "invalid",
				name:  "testInvalid",
				value: "789",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			if err := m.UpdateMetric(tt.args.mtype, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_GetMetric(t *testing.T) {
	// Define some sample data for the test cases
	gaugeMetrics := map[string]interface{}{
		"cpu_usage":    0.75,
		"memory_usage": 0.85,
	}
	counterMetrics := map[string]interface{}{
		"requests_count": 100,
		"errors_count":   5,
	}

	// Populate the test cases
	tests := []struct {
		name    string
		metrics map[MetricType]map[string]interface{}
		mtype   MetricType
		want    map[string]interface{}
	}{
		{
			name:    "EmptyStorage",
			metrics: map[MetricType]map[string]interface{}{},
			mtype:   Gauge,
			want:    map[string]interface{}{},
		},
		{
			name: "GaugeMetrics",
			metrics: map[MetricType]map[string]interface{}{
				Gauge: gaugeMetrics,
			},
			mtype: Gauge,
			want:  gaugeMetrics,
		},
		{
			name: "CounterMetrics",
			metrics: map[MetricType]map[string]interface{}{
				Counter: counterMetrics,
			},
			mtype: Counter,
			want:  counterMetrics,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.metrics,
			}
			got := m.GetMetric(tt.mtype)
			if len(got) == 0 && len(tt.want) == 0 {
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetMetrics(t *testing.T) {
	type fields struct {
		metrics map[MetricType]map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[MetricType]map[string]interface{}
	}{
		{
			name: "Empty Metrics Storage",
			fields: fields{
				metrics: map[MetricType]map[string]interface{}{},
			},
			want: map[MetricType]map[string]interface{}{},
		},
		{
			name: "Metrics Storage with Only Gauge Type Metrics",
			fields: fields{
				metrics: map[MetricType]map[string]interface{}{
					Gauge: {
						"metric1": 10.5,
						"metric2": 20.7,
					},
				},
			},
			want: map[MetricType]map[string]interface{}{
				Gauge: {
					"metric1": 10.5,
					"metric2": 20.7,
				},
			},
		},
		{
			name: "Metrics Storage with Only Counter Type Metrics",
			fields: fields{
				metrics: map[MetricType]map[string]interface{}{
					Counter: {
						"metric1": 5,
						"metric2": 8,
					},
				},
			},
			want: map[MetricType]map[string]interface{}{
				Counter: {
					"metric1": 5,
					"metric2": 8,
				},
			},
		},
		{
			name: "Metrics Storage with Both Gauge and Counter Type Metrics",
			fields: fields{
				metrics: map[MetricType]map[string]interface{}{
					Gauge: {
						"metric1": 10.5,
						"metric2": 20.7,
					},
					Counter: {
						"metric3": 5,
						"metric4": 8,
					},
				},
			},
			want: map[MetricType]map[string]interface{}{
				Gauge: {
					"metric1": 10.5,
					"metric2": 20.7,
				},
				Counter: {
					"metric3": 5,
					"metric4": 8,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			if got := m.GetMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMetricTypeByCode(t *testing.T) {
	type args struct {
		mtype string
	}
	tests := []struct {
		name    string
		args    args
		want    MetricType
		wantErr bool
	}{
		{
			name:    "GaugeMetricType",
			args:    args{mtype: "gauge"},
			want:    Gauge,
			wantErr: false,
		},
		{
			name:    "CounterMetricType",
			args:    args{mtype: "counter"},
			want:    Counter,
			wantErr: false,
		},
		{
			name:    "InvalidMetricType",
			args:    args{mtype: "invalid"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMetricTypeByCode(tt.args.mtype)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetricTypeByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMetricTypeByCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_UpdateJsonMetric(t *testing.T) {
	type fields struct {
		metrics map[MetricType]map[string]interface{}
	}
	type args struct {
		metric *models.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			if err := m.UpdateJsonMetric(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.UpdateJsonMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
