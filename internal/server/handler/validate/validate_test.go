package validate

import (
	"testing"
)

func TestCheckMetricTypeAndName(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		metricName string
		wantErr    bool
	}{
		{
			name:       "Valid parameters",
			metricType: "cpu",
			metricName: "usage",
			wantErr:    false,
		},
		{
			name:       "Missing metric type",
			metricType: "",
			metricName: "usage",
			wantErr:    true,
		},
		{
			name:       "Missing metric name",
			metricType: "cpu",
			metricName: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckMetricTypeAndName(tt.metricType, tt.metricName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckMetricTypeAndName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
