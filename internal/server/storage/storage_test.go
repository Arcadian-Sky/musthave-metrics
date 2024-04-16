package storage

import "testing"

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
				t.Errorf("storage.GetMetricTypeByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("storage.GetMetricTypeByCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
