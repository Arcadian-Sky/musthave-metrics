package sender

import (
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/generated/mocks"
	pb "github.com/Arcadian-Sky/musthave-metrics/internal/agent/generated/protoagent"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendMetricJSONbyGRPC_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAgentServiceClient(ctrl)

	// Устанавливаем ожидаемый вызов
	mockClient.EXPECT().
		SendMetricJSON(gomock.Any(), &pb.MetricJSONRequest{
			JsonString: `{"delta":"10","id":"metric2","type":"counter","value":0}`,
		}).
		Return(&pb.MetricResponse{}, nil).
		Times(1)

	sender := &Sender{
		tcpEndpoint: "localhost:50051",
		cryptoKey:   nil,
		getHash:     "",
		tcpEnabled:  true,
		tcpClient:   mockClient,
	}

	err := sender.SendMetricJSONbyGRPC(map[string]interface{}{
		"id":    "metric2",
		"type":  "counter",
		"delta": "10",
		"value": 0,
	}, "/metric/update")

	assert.NoError(t, err)

	// Устанавливаем ожидаемый вызов
	mockClient.EXPECT().
		SendMetricJSON(gomock.Any(), &pb.MetricJSONRequest{
			JsonString: `{"delta":0,"id":"metric3","type":"gauge","value":10}`,
		}).
		Return(&pb.MetricResponse{}, nil).
		Times(1)

	sender = &Sender{
		tcpEndpoint: "localhost:50051",
		cryptoKey:   nil,
		getHash:     "",
		tcpEnabled:  true,
		tcpClient:   mockClient,
	}

	err = sender.SendMetricJSONbyGRPC(map[string]interface{}{
		"id":    "metric3",
		"type":  "gauge",
		"delta": 0,
		"value": 10,
	}, "/metric/update")

	assert.NoError(t, err)
}

func TestSendValueByGRPC_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAgentServiceClient(ctrl)

	mockClient.EXPECT().
		SendMetric(gomock.Any(), &pb.MetricRequest{
			Type:  "counter",
			Name:  "metric2",
			Value: "100",
		}).
		Return(&pb.MetricResponse{}, nil).
		Times(1)

	sender := &Sender{
		tcpEndpoint: "localhost:50051",
		cryptoKey:   nil,
		getHash:     "",
		tcpEnabled:  true,
		tcpClient:   mockClient,
	}

	err := sender.SendValueByGRPC("counter", "metric2", "100")

	assert.NoError(t, err)
	// Устанавливаем ожидаемый вызов
	mockClient.EXPECT().
		SendMetric(gomock.Any(), &pb.MetricRequest{
			Type:  "gauge",
			Name:  "metric3",
			Value: "10",
		}).
		Return(&pb.MetricResponse{}, nil).
		Times(1)

	sender = &Sender{
		tcpEndpoint: "localhost:50051",
		cryptoKey:   nil,
		getHash:     "",
		tcpEnabled:  true,
		tcpClient:   mockClient,
	}

	err = sender.SendValueByGRPC("gauge", "metric3", "10")

	assert.NoError(t, err)
}
