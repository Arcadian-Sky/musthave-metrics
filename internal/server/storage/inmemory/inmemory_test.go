package inmemory

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/utils"
)

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
			ctx := context.Background()
			got := m.GetMetric(ctx, tt.mtype)
			if len(got) == 0 && len(tt.want) == 0 {
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			ctx := context.Background()
			if err := m.UpdateMetric(ctx, tt.args.mtype, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_GetJSONMetric(t *testing.T) {
	type fields struct {
		metrics map[storage.MetricType]map[string]interface{}
	}
	mtypeGauge := "gauge"
	mtypeCounter := "counter"
	type args struct {
		ctx    context.Context
		metric *models.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty metrics storage",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{},
			},
			args: args{
				ctx: context.Background(),
				metric: &models.Metrics{
					MType: mtypeGauge,
				},
			},
			want:    `{"id":"","type":"gauge","value":0}`,
			wantErr: false,
		},
		{
			name: "metrics storage with only gauge metrics",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": float64(10.5),
						"metric2": float64(20.7),
					},
					storage.Counter: {
						"metric3": int64(5),
						"metric4": int64(8),
					},
				},
			},
			args: args{
				ctx: context.Background(),
				metric: &models.Metrics{
					MType: mtypeGauge,
				},
			},
			want:    `{"id":"","type":"gauge","value":0}`,
			wantErr: false,
		},
		{
			name: "metrics storage with only counter metrics",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": float64(10.5),
						"metric2": float64(20.7),
					},
					storage.Counter: {
						"metric3": int64(5),
						"metric4": int64(8),
					},
				},
			},
			args: args{
				ctx: context.Background(),
				metric: &models.Metrics{
					MType: mtypeCounter,
				},
			},
			want:    `{"id":"","type":"counter","delta":0}`,
			wantErr: false,
		},
		{
			name: "metrics storage with both gauge and counter metrics",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": float64(10.5),
						"metric2": float64(20.7),
					},
					storage.Counter: {
						"metric3": int64(5),
						"metric4": int64(8),
					},
				},
			},
			args: args{
				ctx: context.Background(),
				metric: &models.Metrics{
					MType: mtypeGauge,
				},
			},
			want:    `{"id":"","type":"gauge","value":0}`,
			wantErr: false,
		},
		{
			name: "metrics storage with both gauge and counter metrics 1",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": float64(10.5),
						"metric2": float64(20.7),
					},
					storage.Counter: {
						"metric3": int64(5),
						"metric4": int64(8),
					},
				},
			},
			args: args{
				ctx: context.Background(),
				metric: &models.Metrics{
					MType: mtypeCounter,
					ID:    "metric4",
				},
			},
			want:    `{"id":"metric4","type":"counter","delta":8}`,
			wantErr: false,
		},
		{
			name: "metrics storage with both gauge and counter metrics 2",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": float64(10.5),
						"metric2": float64(20.7),
					},
					storage.Counter: {
						"metric3": int64(5),
						"metric4": int64(8),
					},
				},
			},
			args: args{
				ctx: context.Background(),
				metric: &models.Metrics{
					MType: mtypeGauge,
					ID:    "metric2",
				},
			},
			want:    `{"id":"metric2","type":"gauge","value":20.7}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			// fmt.Printf("before: %v\n", tt.args.metric)
			err := m.GetJSONMetric(tt.args.ctx, tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetJSONMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// fmt.Printf("got: %v\n", tt.args.metric)

			// Преобразуем указатель на структуру в JSON строку
			jsonData, err := json.Marshal(tt.args.metric)
			if err != nil {
				fmt.Println("Ошибка при преобразовании в JSON:", err)
				return
			}

			// fmt.Println(string(jsonData))

			assert.Equal(t, tt.want, string(jsonData))

		})
	}
}

func TestMemStorage_UpdateJSONMetric(t *testing.T) {
	mtypeGauge := "gauge"
	mtypeCounter := "counter"
	type fields struct {
		metrics map[storage.MetricType]map[string]interface{}
	}
	var f100 = float64(100)
	var f200 = float64(200)
	var i300 = int64(300)
	tests := []struct {
		name    string
		fields  fields
		args    models.Metrics
		wantErr bool
		want    string
	}{
		{
			name: "UpdateGaugeMetric",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": float64(10.5),
						"metric2": float64(20.7),
					},
					storage.Counter: {
						"metric3": int64(5),
						"metric4": int64(8),
					},
				},
			},
			args: models.Metrics{
				MType: mtypeGauge,
				ID:    "metric2",
				Value: &f200,
			},
			want:    `200`,
			wantErr: false,
		},
		{
			name: "UpdateCounterMetric",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": float64(10.5),
						"metric2": float64(20.7),
					},
					storage.Counter: {
						"metric3": int64(5),
						"metric4": int64(8),
					},
				},
			},
			args: models.Metrics{
				MType: mtypeCounter,
				ID:    "metric4",
				Delta: &i300,
			},
			want:    `308`,
			wantErr: false,
		},
		{
			name: "InvalidMetricType",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge:   {"metric1": 123.45, "metric2": 678.90},
					storage.Counter: {"metric3": 100, "metric4": 200},
				},
			},
			args: models.Metrics{
				MType: "incorrect",
				ID:    "metric5",
				Value: &f100,
			},
			want:    `null`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			ctx := context.Background()
			// fmt.Printf("m.GetMetrics(ctx): %v\n", m.GetMetrics(ctx))
			err := m.UpdateJSONMetric(ctx, &tt.args)
			// fmt.Printf("m.GetMetrics(ctx)2: %v\n", m.GetMetrics(ctx))

			if (err != nil) && !tt.wantErr {
				t.Errorf("MemStorage.UpdateJSONMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Преобразуем указатель на структуру в JSON строку
			// jsonData, err := json.Marshal(tt.args)
			// if err != nil {
			// 	fmt.Println("Ошибка при преобразовании в JSON:", err)
			// 	return
			// }

			if tt.want != "" {
				mtype, _ := utils.GetMetricTypeByCode(tt.args.MType)
				jsonData, err := json.Marshal(m.metrics[mtype][tt.args.ID])
				if err != nil {
					fmt.Println("Ошибка при преобразовании в JSON:", err)
					return
				}
				fmt.Printf("tt.args: %v\n", string(jsonData))

				assert.Equal(t, tt.want, string(jsonData))
			}

			// fmt.Printf("tt.args: %v\n", tt.args)

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
			ctx := context.Background()
			if got := m.GetMetrics(ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_SetMetrics(t *testing.T) {
	type fields struct {
		metrics map[storage.MetricType]map[string]interface{}
	}
	type args struct {
		ctx     context.Context
		metrics map[storage.MetricType]map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set metrics",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{},
			},
			args: args{
				ctx: context.Background(),
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge: {
						"metric1": 10.5,
						"metric2": 20.7,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			m.SetMetrics(tt.args.ctx, tt.args.metrics)
			assert.Equal(t, tt.args.metrics, m.metrics)
		})
	}
}

func TestMemStorage_UpdateJSONMetrics(t *testing.T) {
	mtypeGauge := "gauge"
	mtypeCounter := "counter"
	type fields struct {
		metrics map[storage.MetricType]map[string]interface{}
	}
	var f100 = float64(100)
	var f200 = float64(200)
	var i300 = int64(300)
	tests := []struct {
		name    string
		fields  fields
		args    []models.Metrics
		wantErr bool
		want    string
	}{
		{
			name: "UpdateGaugeMetric",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge:   {"metric1": 123.45, "metric2": 678.90},
					storage.Counter: {"metric3": 100, "metric4": 200},
				},
			},
			args: []models.Metrics{
				{
					MType: mtypeGauge,
					ID:    "metric4",
					Value: &f200,
				},
			},
			want:    `[{"id":"metric4","type":"gauge","value":200}]`,
			wantErr: false,
		},
		{
			name: "UpdateCounterMetric",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge:   {"metric1": 123.45, "metric2": 678.90},
					storage.Counter: {"metric3": 100, "metric4": 200},
				},
			},
			args: []models.Metrics{
				{
					MType: mtypeCounter,
					ID:    "metric4",
					Delta: &i300,
				},
			},
			want:    `[{"id":"metric4","type":"counter","delta":300}]`,
			wantErr: false,
		},
		{
			name: "InvalidMetricType",
			fields: fields{
				metrics: map[storage.MetricType]map[string]interface{}{
					storage.Gauge:   {"metric1": 123.45, "metric2": 678.90},
					storage.Counter: {"metric3": 100, "metric4": 200},
				},
			},
			args: []models.Metrics{
				{
					MType: "incorrect",
					ID:    "metric5",
					Value: &f100,
				},
			},
			want:    `[{"id":"metric5","type":"incorrect","value":100}]`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			ctx := context.Background()
			err := m.UpdateJSONMetrics(ctx, &tt.args)
			if err != nil && !tt.wantErr {
				t.Errorf("MemStorage.UpdateJSONMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Преобразуем указатель на структуру в JSON строку
			jsonData, err := json.Marshal(tt.args)
			if err != nil {
				fmt.Println("Ошибка при преобразовании в JSON:", err)
				return
			}

			if tt.want != "" {
				assert.Equal(t, tt.want, string(jsonData))
			}

			fmt.Printf("tt.args: %v\n", string(jsonData))
			// fmt.Printf("tt.args: %v\n", tt.args)

		})
	}
}

func TestMemStorage_CreateMemento(t *testing.T) {
	var expected = `{"metrics":{"counter":{"metric1":100,"metric2":200},"gauge":{"metric3":300.5}}}`
	var metrics = make(map[storage.MetricType]map[string]interface{})
	metrics[storage.Counter] = map[string]interface{}{
		"metric1": 100,
		"metric2": 200,
	}
	metrics[storage.Gauge] = map[string]interface{}{
		"metric3": 300.5,
	}
	mStorage := NewMemStorage()
	ctx := context.Background()
	mStorage.SetMetrics(ctx, metrics)

	// Вызываем функцию CreateMemento
	memento := mStorage.CreateMemento()

	// // fmt.Printf("response.Body: %v\n", response.Body)
	respJSONData, err := json.Marshal(memento)
	if err != nil {
		t.Errorf("Failed to marshal empty metric: %v", err)
	}

	// Сравниваем результат с ожидаемым
	if !reflect.DeepEqual(string(respJSONData), expected) {
		t.Errorf("Expected result %v, got %v", string(respJSONData), expected)
	}
}

func TestMemStorage_RestoreFromMemento(t *testing.T) {
	var expected = `{"counter":{"metric20":200,"metric5":100},"gauge":{"metric100":300.5}}`
	var metrics = make(map[storage.MetricType]map[string]interface{})
	metrics[storage.Counter] = map[string]interface{}{
		"metric5":  100,
		"metric20": 200,
	}
	metrics[storage.Gauge] = map[string]interface{}{
		"metric100": 300.5,
	}

	mStorage := NewMemStorage()
	ctx := context.Background()
	mStorage.SetMetrics(ctx, metrics)

	s := &storage.Memento{}
	s.SetMetrics(metrics)

	mStorage.RestoreFromMemento(s)

	// // fmt.Printf("response.Body: %v\n", response.Body)
	respJSONData, err := json.Marshal(mStorage.metrics)
	if err != nil {
		t.Errorf("Failed to marshal empty metric: %v", err)
	}

	// Сравниваем результат с ожидаемым
	if !reflect.DeepEqual(string(respJSONData), expected) {
		t.Errorf("Expected result %v, got %v", string(respJSONData), expected)
	}
}
