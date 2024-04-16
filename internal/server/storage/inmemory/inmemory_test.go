package inmemory

import (
	"context"
	"reflect"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

func TestMemStorage_UpdateMetric(t *testing.T) {
	type fields struct {
		metrics map[storage.MetricType]map[string]interface{}
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
				metrics: map[storage.MetricType]map[string]interface{}{},
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
				metrics: map[storage.MetricType]map[string]interface{}{},
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
				metrics: map[storage.MetricType]map[string]interface{}{},
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
			ctx := context.TODO()
			if err := m.UpdateMetric(ctx, tt.args.mtype, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
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
		metrics map[storage.MetricType]map[string]interface{}
		mtype   storage.MetricType
		want    map[string]interface{}
	}{
		{
			name:    "EmptyStorage",
			metrics: map[storage.MetricType]map[string]interface{}{},
			mtype:   storage.Gauge,
			want:    map[string]interface{}{},
		},
		{
			name: "GaugeMetrics",
			metrics: map[storage.MetricType]map[string]interface{}{
				storage.Gauge: gaugeMetrics,
			},
			mtype: storage.Gauge,
			want:  gaugeMetrics,
		},
		{
			name: "CounterMetrics",
			metrics: map[storage.MetricType]map[string]interface{}{
				storage.Counter: counterMetrics,
			},
			mtype: storage.Counter,
			want:  counterMetrics,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.metrics,
			}
			ctx := context.TODO()
			got := m.GetMetric(ctx, tt.mtype)
			if len(got) == 0 && len(tt.want) == 0 {
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetMetrics(t *testing.T) {
	type fields struct {
		metrics map[storage.MetricType]map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[storage.MetricType]map[string]interface{}
	}{
		{
			name: "Empty Metrics Storage",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{},
			},
			want: map[storage.MetricType]map[string]interface{}{},
		},
		{
			name: "Metrics Storage with Only storage.Gauge Type Metrics",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": 10.5,
						"metric2": 20.7,
					},
				},
			},
			want: map[storage.MetricType]map[string]interface{}{
				storage.Gauge: {
					"metric1": 10.5,
					"metric2": 20.7,
				},
			},
		},
		{
			name: "Metrics Storage with Only storage.Counter Type Metrics",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Counter: {
						"metric1": 5,
						"metric2": 8,
					},
				},
			},
			want: map[storage.MetricType]map[string]interface{}{
				storage.Counter: {
					"metric1": 5,
					"metric2": 8,
				},
			},
		},
		{
			name: "Metrics Storage with Both storage.Gauge and storage.Counter Type Metrics",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": 10.5,
						"metric2": 20.7,
					},
					storage.Counter: {
						"metric3": 5,
						"metric4": 8,
					},
				},
			},
			want: map[storage.MetricType]map[string]interface{}{
				storage.Gauge: {
					"metric1": 10.5,
					"metric2": 20.7,
				},
				storage.Counter: {
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
			ctx := context.TODO()
			if got := m.GetMetrics(ctx); !reflect.DeepEqual(got, tt.want) {
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
		want    storage.MetricType
		wantErr bool
	}{
		{
			name:    "GaugeMetricType",
			args:    args{mtype: "gauge"},
			want:    storage.Gauge,
			wantErr: false,
		},
		{
			name:    "CounterMetricType",
			args:    args{mtype: "counter"},
			want:    storage.Counter,
			wantErr: false,
		},
		// {
		// 	name:    "InvalidMetricType",
		// 	args:    args{mtype: "invalid"},
		// 	want:    0,
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.GetMetricTypeByCode(tt.args.mtype)
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.GetMetricTypeByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("storage.GetMetricTypeByCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_UpdateJSONMetric(t *testing.T) {
	type fields struct {
		metrics map[storage.MetricType]map[string]interface{}
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
			ctx := context.TODO()
			if err := m.UpdateJSONMetric(ctx, tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.UpdateJSONMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
