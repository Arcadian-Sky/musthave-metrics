package controller

import (
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
)

func TestNewCollectAndSendMetricsService(t *testing.T) {

	type fields struct {
		config flags.Config
	}
	type args struct {
		mType  string
		mName  string
		mValue interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// {
		// 	name: "Test error sending metric value",
		// 	fields: fields{
		// 		config: *flags.SetDefault(),
		// 	},
		// 	args: args{
		// 		mType:  "gauge",
		// 		mName:  "testMetric",
		// 		mValue: 10,
		// 	},
		// 	wantErr: true,
		// },
		{
			name: "Test error sending metric value",
			fields: fields{
				config: flags.Config{},
			},
			args: args{
				mType:  "gauge",
				mName:  "testMetric",
				mValue: 10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CollectAndSendMetricsService{
				config: tt.fields.config,
			}
			if err := c.sendMetricValue(tt.args.mType, tt.args.mName, tt.args.mValue); (err != nil) != tt.wantErr {
				t.Errorf("CollectAndSendMetricsService.sendMetricValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollectAndSendMetricsService_send(t *testing.T) {
	// metrics := make(map[string]interface{})
	// pollCount := 10 // Пример значения pollCount

	// service := NewCollectAndSendMetricsService(*flags.SetDefault())
	// err := service.send(metrics, pollCount)
	// assert.Error(t, err, "Функция должна вернуть errоr")
}

func TestCollectAndSendMetricsService_Run(t *testing.T) {
	type fields struct {
		config flags.Config
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CollectAndSendMetricsService{
				config: tt.fields.config,
			}
			c.Run()
		})
	}
}

func TestCollectAndSendMetricsService_sendMetricValue(t *testing.T) {
	type fields struct {
		config flags.Config
	}
	type args struct {
		mType  string
		mName  string
		mValue interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// {
		// 	name: "Test with valid data",
		// 	fields: fields{
		// 		config: *flags.SetDefault(),
		// 	},
		// 	args: args{
		// 		mType:  "gauge",
		// 		mName:  "testMetric",
		// 		mValue: 123,
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CollectAndSendMetricsService{
				config: tt.fields.config,
			}
			if err := c.sendMetricValue(tt.args.mType, tt.args.mName, tt.args.mValue); (err != nil) != tt.wantErr {
				t.Errorf("CollectAndSendMetricsService.sendMetricValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
