package storage

import "testing"

func TestMemStorage_SaveMetrics(t *testing.T) {
	type fields struct {
		metrics map[MetricType]map[string]interface{}
	}
	type args struct {
		config Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				metrics: tt.fields.metrics,
			}
			m.SaveMetrics(tt.args.config)
		})
	}
}