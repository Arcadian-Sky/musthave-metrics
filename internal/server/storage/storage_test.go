package storage

import (
	"reflect"
	"testing"
	"time"
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
		// TODO: Add test cases.
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
	type fields struct {
		metrics map[MetricType]map[string]interface{}
	}
	type args struct {
		mtype MetricType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			if got := m.GetMetric(tt.args.mtype); !reflect.DeepEqual(got, tt.want) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		time.Sleep(2 * time.Second) // Sleep for 2 seconds
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
