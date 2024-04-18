package controller

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/models"
	"github.com/stretchr/testify/assert"
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
	metrics := make(map[string]interface{})
	pollCount := 10 // Пример значения pollCount

	service := NewCollectAndSendMetricsService(*flags.SetDefault())
	err := service.send(metrics, pollCount)
	assert.Error(t, err, "Функция должна вернуть errоr")
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
		{
			name: "Test with valid data",
			fields: fields{
				config: *flags.SetDefault(),
			},
			args: args{
				mType:  "gauge",
				mName:  "testMetric",
				mValue: 123,
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

// func TestSendPack(t *testing.T) {
// 	c := &CollectAndSendMetricsService{} // Создаем экземпляр сервиса

// 	// Подготовка данных для теста
// 	metrics := map[string]interface{}{
// 		"metric1": 100.0,
// 		"metric2": 200.0,
// 	}
// 	pollCount := 5

// 	// Вызываем метод, который тестируем
// 	err := c.sendPack(metrics, pollCount)

// 	// Проверяем, что ошибки нет
// 	assert.NoError(t, err)
// }

func TestSendMetricJSONValues(t *testing.T) {
	// Создаем фейковый HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем URL запроса
		if r.URL.Path != "/updates" {
			t.Errorf("Expected URL path '/updates', got '%s'", r.URL.Path)
		}
		// Проверяем Content-Type заголовок
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}
		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		// Проверяем тело запроса
		expectedBody := `[{"id":"metric1","type":"gauge","value":100},{"id":"metric2","type":"gauge","value":200},{"id":"PollCount","type":"counter","delta":5}]`
		if string(body) != expectedBody {
			t.Errorf("Unexpected request body. Expected '%s', got '%s'", expectedBody, string(body))
		}
		// Отправляем успешный ответ
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	fmt.Printf("server.URL: %v\n", server.URL)

	config, _ := flags.Parse()
	config.SetConfigServer(server.URL)

	// Создаем экземпляр CollectAndSendMetricsService с фейковым сервером
	c := &CollectAndSendMetricsService{
		config: config,
	}
	var fl100 = float64(100)
	var fl200 = float64(200)
	var de5 = int64(5)
	// // Тестируемый вызов
	err := c.sendMetricJSONValues([]interface{}{
		models.Metrics{ID: "metric1", MType: "gauge", Value: &fl100},
		models.Metrics{ID: "metric2", MType: "gauge", Value: &fl200},
		models.Metrics{ID: "PollCount", MType: "counter", Delta: &de5},
	})
	if err != nil {
		t.Errorf("sendMetricJSONValues failed: %v", err)
	}
}
