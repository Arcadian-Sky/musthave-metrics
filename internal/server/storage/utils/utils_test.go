package utils

import (
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

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
			got, err := GetMetricTypeByCode(tt.args.mtype)
			if (err != nil) != tt.wantErr {
				t.Errorf("utils.GetMetricTypeByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("utils.GetMetricTypeByCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
