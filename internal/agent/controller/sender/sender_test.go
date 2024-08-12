package sender

import (
	"fmt"
	"reflect"
	"testing"

	mockpb "github.com/Arcadian-Sky/musthave-metrics/gen/mocks/api/metrics/v1"
	pb "github.com/Arcadian-Sky/musthave-metrics/gen/proto/api/metrics/v1"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/models"
	"github.com/golang/mock/gomock"
)

// func TestSendMetricJSONbyGRPC_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockClient := mocks.NewMockMetricsServiceClient(ctrl)

// 	// Устанавливаем ожидаемый вызов
// 	mockClient.EXPECT().
// 		SendMetricJSON(gomock.Any(), &pb.MetricJSONRequest{
// 			JsonString: `{"delta":"10","id":"metric2","type":"counter","value":0}`,
// 		}).
// 		Return(&pb.MetricResponse{}, nil).
// 		Times(1)

// 	sender := &Sender{
// 		tcpEndpoint: "localhost:50051",
// 		cryptoKey:   nil,
// 		getHash:     "",
// 		tcpEnabled:  true,
// 		tcpClient:   mockClient,
// 	}

// 	err := sender.SendMetricJSONbyGRPC(map[string]interface{}{
// 		"id":    "metric2",
// 		"type":  "counter",
// 		"delta": "10",
// 		"value": 0,
// 	}, "/metric/update")

// 	assert.NoError(t, err)

// 	// Устанавливаем ожидаемый вызов
// 	mockClient.EXPECT().
// 		SendMetricJSON(gomock.Any(), &pb.MetricJSONRequest{
// 			JsonString: `{"delta":0,"id":"metric3","type":"gauge","value":10}`,
// 		}).
// 		Return(&pb.MetricResponse{}, nil).
// 		Times(1)

// 	sender = &Sender{
// 		tcpEndpoint: "localhost:50051",
// 		cryptoKey:   nil,
// 		getHash:     "",
// 		tcpEnabled:  true,
// 		tcpClient:   mockClient,
// 	}

// 	err = sender.SendMetricJSONbyGRPC(map[string]interface{}{
// 		"id":    "metric3",
// 		"type":  "gauge",
// 		"delta": 0,
// 		"value": 10,
// 	}, "/metric/update")

// 	assert.NoError(t, err)
// }

// func TestSendValueByGRPC_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockClient := mocks.NewMockAgentServiceClient(ctrl)

// 	mockClient.EXPECT().
// 		SendMetric(gomock.Any(), &pb.MetricRequest{
// 			Type:  "counter",
// 			Name:  "metric2",
// 			Value: "100",
// 		}).
// 		Return(&pb.MetricResponse{}, nil).
// 		Times(1)

// 	sender := &Sender{
// 		tcpEndpoint: "localhost:50051",
// 		cryptoKey:   nil,
// 		getHash:     "",
// 		tcpEnabled:  true,
// 		tcpClient:   mockClient,
// 	}

// 	err := sender.SendValueByGRPC("counter", "metric2", "100")

// 	assert.NoError(t, err)
// 	// Устанавливаем ожидаемый вызов
// 	mockClient.EXPECT().
// 		SendMetric(gomock.Any(), &pb.MetricRequest{
// 			Type:  "gauge",
// 			Name:  "metric3",
// 			Value: "10",
// 		}).
// 		Return(&pb.MetricResponse{}, nil).
// 		Times(1)

// 	sender = &Sender{
// 		tcpEndpoint: "localhost:50051",
// 		cryptoKey:   nil,
// 		getHash:     "",
// 		tcpEnabled:  true,
// 		tcpClient:   mockClient,
// 	}

// 	err = sender.SendValueByGRPC("gauge", "metric3", "10")

// 	assert.NoError(t, err)
// }

func Test_convertModelMetricsToRPCMetrics(t *testing.T) {
	type args struct {
		metric *models.Metrics
	}
	tests := []struct {
		name string
		args args
		want *pb.Metric
	}{
		{
			name: "Counter metric with Delta",
			args: args{
				metric: &models.Metrics{
					ID:    "metric1",
					MType: "counter",
					Delta: int64Ptr(10),
				},
			},
			want: &pb.Metric{
				Id:   "metric1",
				Type: pb.Type_TYPE_COUNTER,
				Mvalue: &pb.Metric_Delta{
					Delta: 10,
				},
			},
		},
		{
			name: "Gauge metric with Value",
			args: args{
				metric: &models.Metrics{
					ID:    "metric2",
					MType: "gauge",
					Value: float64Ptr(3.14),
				},
			},
			want: &pb.Metric{
				Id:   "metric2",
				Type: pb.Type_TYPE_GAUGE,
				Mvalue: &pb.Metric_Value{
					Value: 3.14,
				},
			},
		},
		{
			name: "Gauge metric with nil Value",
			args: args{
				metric: &models.Metrics{
					ID:    "metric4",
					MType: "gauge",
					Value: nil,
				},
			},
			want: &pb.Metric{
				Id:     "metric4",
				Type:   pb.Type_TYPE_GAUGE,
				Mvalue: nil,
			},
		},
		{
			name: "Counter metric with nil Delta",
			args: args{
				metric: &models.Metrics{
					ID:    "metric5",
					MType: "counter",
					Delta: nil,
				},
			},
			want: &pb.Metric{
				Id:     "metric5",
				Type:   pb.Type_TYPE_COUNTER,
				Mvalue: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertModelMetricsToRPCMetrics(tt.args.metric); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertModelMetricsToRPCMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions to create pointers for test values
func int64Ptr(value int64) *int64 {
	return &value
}

func float64Ptr(value float64) *float64 {
	return &value
}

func Test_convertStringToMetricType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want pb.Type
	}{
		{
			name: "Gauge type",
			args: args{
				s: "gauge",
			},
			want: pb.Type_TYPE_GAUGE,
		},
		{
			name: "Counter type",
			args: args{
				s: "counter",
			},
			want: pb.Type_TYPE_COUNTER,
		},
		{
			name: "Unknown type",
			args: args{
				s: "unknown",
			},
			want: pb.Type_TYPE_UNSPECIFIED,
		},
		{
			name: "Empty string",
			args: args{
				s: "",
			},
			want: pb.Type_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertStringToMetricType(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertStringToMetricType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSender_SendValueByGRPC(t *testing.T) {
	type args struct {
		mType  string
		mName  string
		mValue string
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockpb.NewMockMetricsServiceClient(ctrl) // Замените на ваш Mock клиент
	sender := &Sender{tcpClient: mockClient}
	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "Successful request with string value",
			args: args{
				mType:  "gauge",
				mName:  "metric1",
				mValue: "123.45",
			},
			setup: func() {
				mockClient.EXPECT().
					UpdateMetric(gomock.Any(), &pb.UpdateMetricRequest{
						Id:    "metric1",
						Type:  pb.Type_TYPE_GAUGE,
						Value: "123.45",
					}).
					Return(&pb.UpdateMetricResponse{}, nil)
			},
			wantErr: false,
		},
		{
			name: "GRPC call returns error",
			args: args{
				mType:  "gauge",
				mName:  "metric3",
				mValue: "678.90",
			},
			setup: func() {
				mockClient.EXPECT().
					UpdateMetric(gomock.Any(), &pb.UpdateMetricRequest{
						Id:    "metric3",
						Type:  pb.Type_TYPE_GAUGE,
						Value: "678.90",
					}).
					Return(nil, fmt.Errorf("GRPC error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := sender.SendValueByGRPC(tt.args.mType, tt.args.mName, tt.args.mValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sender.SendValueByGRPC() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSender_SendMetricJSONbyGRPC(t *testing.T) {
	type args struct {
		m      any
		method string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockpb.NewMockMetricsServiceClient(ctrl)                        // Замените на ваш Mock клиент
	sender := &Sender{tcpClient: mockClient, getHash: "testhash", cryptoKey: nil} // Настройте Sender как необходимо

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "Send single metric successfully",
			args: args{
				m: models.Metrics{
					ID:    "metric1",
					MType: "gauge",
					Value: float64Pointer(12.34),
				},
				method: UpdatePathOne,
			},
			setup: func() {
				mockClient.EXPECT().
					UpdateJSONMetric(gomock.Any(), &pb.UpdateJSONMetricRequest{
						Metric: &pb.Metric{
							Id:   "metric1",
							Type: pb.Type_TYPE_GAUGE,
							Mvalue: &pb.Metric_Value{
								Value: 12.34,
							},
						},
					}).
					Return(&pb.UpdateJSONMetricResponse{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Send multiple metrics successfully",
			args: args{
				m: []models.Metrics{
					{
						ID:    "metric2",
						MType: "counter",
						Delta: int64Pointer(42),
					},
					{
						ID:    "metric3",
						MType: "gauge",
						Value: float64Pointer(56.78),
					},
				},
				method: UpdatePathPack,
			},
			setup: func() {
				mockClient.EXPECT().
					UpdateJSONMetrics(gomock.Any(), &pb.UpdateJSONMetricsRequest{
						Metrics: []*pb.Metric{
							{
								Id:   "metric2",
								Type: pb.Type_TYPE_COUNTER,
								Mvalue: &pb.Metric_Delta{
									Delta: 42,
								},
							},
							{
								Id:   "metric3",
								Type: pb.Type_TYPE_GAUGE,
								Mvalue: &pb.Metric_Value{
									Value: 56.78,
								},
							},
						},
					}).
					Return(&pb.UpdateJSONMetricsResponse{}, nil)
			},
			wantErr: false,
		},
		{
			// Некорректный тип
			name: "Unsupported type",
			args: args{
				m:      123,
				method: UpdatePathOne,
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "Unsupported method",
			args: args{
				m:      models.Metrics{},
				method: "unsupportedMethod",
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "GRPC error",
			args: args{
				m: models.Metrics{
					ID:    "metric4",
					MType: "gauge",
					Value: float64Pointer(78.90),
				},
				method: UpdatePathOne,
			},
			setup: func() {
				mockClient.EXPECT().
					UpdateJSONMetric(gomock.Any(), &pb.UpdateJSONMetricRequest{
						Metric: &pb.Metric{
							Id:   "metric4",
							Type: pb.Type_TYPE_GAUGE,
							Mvalue: &pb.Metric_Value{
								Value: 78.90,
							},
						},
					}).
					Return(nil, fmt.Errorf("GRPC error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := sender.SendMetricJSONbyGRPC(tt.args.m, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sender.SendMetricJSONbyGRPC() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Помощь для указателя на float64
func float64Pointer(f float64) *float64 {
	return &f
}

// Помощь для указателя на int64
func int64Pointer(i int64) *int64 {
	return &i
}
