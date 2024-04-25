package controller

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/controller/sender"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/models"
	"github.com/stretchr/testify/assert"
)

func TestNewCollectAndSendMetricsService(t *testing.T) {

	type fields struct {
		config flags.Config
		sender *sender.Sender
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
				sender: sender.NewSender(&flags.Config{}),
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
				sender: tt.fields.sender,
			}
			if err := c.sender.SendMetricValue(tt.args.mType, tt.args.mName, tt.args.mValue); (err != nil) != tt.wantErr {
				t.Errorf("CollectAndSendMetricsService.sendMetricValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollectAndSendMetricsService_send(t *testing.T) {
	metrics := make(map[string]interface{})
	pollCount := int64(10) // Пример значения pollCount

	service := NewCollectAndSendMetricsService(flags.SetDefault())
	err := service.send(metrics, pollCount)
	assert.Error(t, err, "Функция должна вернуть errоr")
}

func TestCollectAndSendMetricsService_sendMetricValue(t *testing.T) {
	type fields struct {
		config flags.Config
		sender *sender.Sender
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
				sender: sender.NewSender(&flags.Config{}),
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
				sender: tt.fields.sender,
			}
			if err := c.sender.SendMetricValue(tt.args.mType, tt.args.mName, tt.args.mValue); (err != nil) != tt.wantErr {
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

func TestMakePack(t *testing.T) {
	// Создаем экземпляр CollectAndSendMetricsService
	c := &CollectAndSendMetricsService{}

	// Создаем тестовые данные для метрик
	metrics := map[string]interface{}{
		"metric1": 100.0,
		"metric2": 200.0,
	}

	// Создаем тестовое значение для pollCount
	pollCount := int64(10)

	// Вызываем метод makePack
	pack := c.makePack(metrics, pollCount)
	// Проверяем, что pack не nil
	if pack == nil {
		t.Error("Expected pack to not be nil")
	}

	// Проверяем, что длина pack равна количеству метрик + 1
	expectedLength := len(metrics) + 1
	if len(pack) != expectedLength {
		t.Errorf("Expected pack length to be %d, but got %d", expectedLength, len(pack))
	}

	// Проверяем, что pack содержит корректные метрики
	for i, metric := range pack {
		if i < len(metrics) {
			// Проверяем метрики типа gauge
			m, ok := metric.(models.Metrics)
			if !ok {
				t.Errorf("Expected metric at index %d to be type models.Metrics", i)
			}
			if m.MType != "gauge" {
				t.Errorf("Expected metric type to be 'gauge', but got '%s'", m.MType)
			}
		} else {
			// Проверяем метрику типа counter
			m, ok := metric.(models.Metrics)
			if !ok {
				t.Errorf("Expected metric at index %d to be type models.Metrics", i)
			}
			if m.MType != "counter" {
				t.Errorf("Expected metric type to be 'counter', but got '%s'", m.MType)
			}
		}
	}
}

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

	conf, _ := flags.Parse()
	conf.SetConfigServer(server.URL)

	// Создаем экземпляр CollectAndSendMetricsService с фейковым сервером
	c := &CollectAndSendMetricsService{
		config: conf,
		sender: sender.NewSender(&conf),
	}
	var fl100 = float64(100)
	var fl200 = float64(200)
	var de5 = int64(5)
	// // Тестируемый вызов
	err := c.sender.SendMetricJSON([]interface{}{
		models.Metrics{ID: "metric1", MType: "gauge", Value: &fl100},
		models.Metrics{ID: "metric2", MType: "gauge", Value: &fl200},
		models.Metrics{ID: "PollCount", MType: "counter", Delta: &de5},
	}, "/updates")
	if err != nil {
		t.Errorf("sendMetricJSONValues failed: %v", err)
	}
}
