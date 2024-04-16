package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// func TestUpdateMetricsHandler(t *testing.T) {
// 	// Создаем фейковое хранилище
// 	storage.Storage = storage.NewMemStorage()

// 	// Создаем HTTP запрос для обновления метрики
// 	req, err := http.NewRequest("POST", "/update/gauge/example_metric/42", nil)
// 	assert.NoError(t, err)

// 	// Создаем ResponseRecorder для записи ответа сервера
// 	w := httptest.NewRecorder()

// 	// Вызыываем обработчик
// 	UpdateMetricsHandlerFunc(w, req)

// 	res := w.Result()

// 	defer res.Body.Close()

// 	// Проверяем код ответа
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Проверяем тело ответа
// 	assert.Contains(t, w.Body.String(), "0: map[example_metric:42]\n")
// }

// func TestUpdateMetricsHandlerBadRequest(t *testing.T) {
// 	// Создаем фейковое хранилище
// 	storage.Storage = storage.NewMemStorage()

// 	// Создаем HTTP запрос с некорректным типом метрики
// 	req, err := http.NewRequest("POST", "/update/invalid_type/example_metric/42", nil)
// 	assert.NoError(t, err)

// 	// Создаем ResponseRecorder для записи ответа сервера
// 	w := httptest.NewRecorder()

// 	// Вызыываем обработчик
// 	UpdateMetricsHandlerFunc(w, req)

// 	res := w.Result()

// 	defer res.Body.Close()

// 	// Проверяем код ответа
// 	assert.Equal(t, http.StatusBadRequest, w.Code)

// 	// Проверяем тело ответа
// 	// assert.Contains(t, rr.Body.String(), "Metric type validation failed")
// }

func InitRouter() chi.Router {

	handler := NewHandler(inmemory.NewMemStorage())

	r := chi.NewRouter()

	r.Head("/", func(rw http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "text/plain")
	})
	// GET http://localhost:8080/value/counter/testSetGet163
	r.Get("/", handler.MetricsHandlerFunc)
	r.Get("/ping", handler.PingDB)
	r.Post("/updates/", handler.UpdateJSONMetricsHandlerFunc)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handler.UpdateJSONMetricHandlerFunc)
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}/", handler.UpdateMetricsHandlerFunc)
			r.Get("/{name}/{value}/", handler.UpdateMetricsHandlerFunc)

		})
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", handler.GetMetricsJSONHandlerFunc)
		r.Get("/", handler.GetMetricHandlerFunc)
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/", handler.GetMetricHandlerFunc)
			r.Get("/{name}", handler.GetMetricHandlerFunc)
			r.Get("/{name}/", handler.GetMetricHandlerFunc)
		})
	})

	return r
}

func TestHandler_GetMetricHandlerFunc(t *testing.T) {
	tests := []struct {
		name          string
		requestPath   string
		expectedType  string
		expectedName  string
		expectedCode  int
		expectedValue string
	}{

		{
			name:          "not valid request 1",
			requestPath:   "/value/gauge/someName/",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  404,
			expectedValue: "",
		},
		{
			name:          "not valid request 2",
			requestPath:   "/value/counter/someName/",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  404,
			expectedValue: "",
		},
		{
			name:          "not valid request 3",
			requestPath:   "/value/error/someName/",
			expectedType:  "",
			expectedName:  "",
			expectedCode:  400,
			expectedValue: "",
		},
		{
			name:          "not valid request count no name",
			requestPath:   "/value/counter/",
			expectedType:  "count",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "metric name not provided",
		},
		{
			name:          "not valid request gauge no name",
			requestPath:   "/value/gauge/",
			expectedType:  "gauge",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "metric name not provided",
		},
		{
			name:          "not valid request name",
			requestPath:   "/value/error/",
			expectedType:  "",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "metric name not provided",
		},
		{
			name:          "not valid request gauge no type",
			requestPath:   "/value/",
			expectedType:  "",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "metric type not provided",
		},
	}

	testServer := httptest.NewServer(InitRouter())
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, testServer.URL+tt.requestPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, response.StatusCode)

			respBody, _ := io.ReadAll(response.Body)
			fmt.Println(string(respBody))
			assert.Contains(t, string(respBody), tt.expectedValue)
			defer response.Body.Close()
		})
	}
}

func TestHandler_UpdateMetricsHandlers(t *testing.T) {
	tests := []struct {
		name          string
		requestPath   string
		expectedType  string
		expectedName  string
		expectedCode  int
		expectedValue string
	}{
		{
			name:          "valid request gauge",
			requestPath:   "/update/gauge/someName/100.001",
			expectedType:  "gauge",
			expectedName:  "someName",
			expectedCode:  200,
			expectedValue: "map[someName:100.001]",
		},
		{
			name:          "valid request count",
			requestPath:   "/update/counter/someName/100",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  200,
			expectedValue: "map[someName:100]",
		},
		{
			name:          "not valid request error counter",
			requestPath:   "/update/counter/someName/100.001",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  400,
			expectedValue: "",
		},
		{
			name:          "not valid request 1",
			requestPath:   "/update/gauge/someName/",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  400,
			expectedValue: "invalid",
		},
		{
			name:          "not valid request 2",
			requestPath:   "/update/counter/someName/",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  400,
			expectedValue: "invalid",
		},
		{
			name:          "not valid request 3",
			requestPath:   "/update/error/someName/",
			expectedType:  "",
			expectedName:  "",
			expectedCode:  400,
			expectedValue: "invalid metric type",
		},
		{
			name:          "not valid request count no name",
			requestPath:   "/update/counter/",
			expectedType:  "count",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "metric name not provided",
		},
		{
			name:          "not valid request error no name 2",
			requestPath:   "/update/error/",
			expectedType:  "count",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "metric name not provided",
		},
		{
			name:          "not valid request gauge no name",
			requestPath:   "/update/gauge/",
			expectedType:  "gauge",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "metric name not provided",
		},
		{
			name:          "not valid request gauge no type",
			requestPath:   "/update/",
			expectedType:  "",
			expectedName:  "",
			expectedCode:  400,
			expectedValue: "Request body is empty",
		},
	}

	testServer := httptest.NewServer(InitRouter())
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, testServer.URL+tt.requestPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, response.StatusCode)

			respBody, _ := io.ReadAll(response.Body)
			fmt.Println(string(respBody))
			assert.Contains(t, string(respBody), tt.expectedValue)
			defer response.Body.Close()
		})
	}
}

func TestHandler_MetricsHandlerFunc(t *testing.T) {
	tests := []struct {
		name          string
		requestPath   string
		expectedType  string
		expectedName  string
		expectedCode  int
		expectedValue string
	}{

		{
			name:          "valid request 1",
			requestPath:   "/",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  200,
			expectedValue: "",
		},
	}

	testServer := httptest.NewServer(InitRouter())
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, testServer.URL+tt.requestPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, response.StatusCode)

			respBody, _ := io.ReadAll(response.Body)
			fmt.Println(string(respBody))
			assert.Contains(t, string(respBody), tt.expectedValue)
			defer response.Body.Close()
		})
	}
}

func TestHandler_UpdateMetricsJSONHandlerFunc(t *testing.T) {
	value := 100.0
	var delta int64 = 100
	// var value0 float64 = 0
	tests := []struct {
		requestPath     string
		name            string
		requestBody     string
		expectedCode    int
		expectedMetrics models.Metrics
	}{
		{
			requestPath:  "/value",
			name:         "valid request gauge",
			requestBody:  `{"id": "metric1", "type": "gauge", "value": 100}`,
			expectedCode: http.StatusOK,
			expectedMetrics: models.Metrics{
				ID:    "metric1",
				MType: "gauge",
				Value: &value,
			},
		},
		{
			requestPath:  "/value",
			name:         "valid request counter",
			requestBody:  `{"id": "metric2", "type": "counter", "delta": 100}`,
			expectedCode: http.StatusOK,
			expectedMetrics: models.Metrics{
				ID:    "metric1",
				MType: "gauge",
				Delta: &delta,
			},
		},
		{
			requestPath:     "/value",
			name:            "empty request",
			requestBody:     `{}`, // Empty request body
			expectedCode:    http.StatusBadRequest,
			expectedMetrics: models.Metrics{},
		},
		{
			requestPath:  "/value",
			name:         "only ID request",
			requestBody:  `{"id": "metric1"}`,
			expectedCode: http.StatusBadRequest,
			expectedMetrics: models.Metrics{
				ID: "metric1",
				// Value: &value0,
			},
		},
		{
			requestPath:  "/value",
			name:         "ID and type request",
			requestBody:  `{"id": "metric1", "type": "gauge"}`,
			expectedCode: http.StatusOK,
			expectedMetrics: models.Metrics{
				ID:    "metric1",
				MType: "gauge",
				// Value: &value0,
			},
		},
	}

	testServer := httptest.NewServer(InitRouter())
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := bytes.NewBufferString(tt.requestBody)
			request, err := http.NewRequest(http.MethodPost, testServer.URL+"/update/", requestBody)
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Content-Type", "application/json")

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, response.StatusCode)

			var result models.Metrics
			if tt.expectedCode == http.StatusOK {
				fmt.Printf("response.Body: %v\n", response.Body)
				err := json.NewDecoder(response.Body).Decode(&result)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Printf("result: %v\n", result)

				var buf bytes.Buffer
				encoder := json.NewEncoder(&buf)
				assert.Equal(t, encoder.Encode(result), encoder.Encode(tt.expectedMetrics))
			}

			defer response.Body.Close()
		})
	}
}

func TestHandler_GetMetricsJSONHandlerFunc(t *testing.T) {
	value := 100.0
	var delta int64 = 100
	// var value0 float64 = 0
	tests := []struct {
		requestPath     string
		name            string
		requestBody     string
		expectedCode    int
		expectedMetrics models.Metrics
	}{
		{
			requestPath:  "/value",
			name:         "valid request",
			requestBody:  `{"id": "metric1", "type": "gauge", "value": 100}`,
			expectedCode: http.StatusOK,
			expectedMetrics: models.Metrics{
				ID:    "metric1",
				MType: "gauge",
				Value: &value,
			},
		},
		{
			requestPath:  "/value",
			name:         "valid request",
			requestBody:  `{"id": "metric2", "type": "counter", "delta": 100}`,
			expectedCode: http.StatusOK,
			expectedMetrics: models.Metrics{
				ID:    "metric1",
				MType: "gauge",
				Delta: &delta,
			},
		},
		{
			requestPath:     "/value",
			name:            "empty request",
			requestBody:     `{}`, // Empty request body
			expectedCode:    http.StatusOK,
			expectedMetrics: models.Metrics{},
		},
		{
			requestPath:  "/value",
			name:         "only ID request",
			requestBody:  `{"id": "metric1"}`,
			expectedCode: http.StatusOK,
			expectedMetrics: models.Metrics{
				ID: "metric1",
				// Value: &value0,
			},
		},
		{
			requestPath:  "/value",
			name:         "ID and type request",
			requestBody:  `{"id": "metric1", "type": "gauge"}`,
			expectedCode: http.StatusOK,
			expectedMetrics: models.Metrics{
				ID:    "metric1",
				MType: "gauge",
				// Value: &value0,
			},
		},
	}

	testServer := httptest.NewServer(InitRouter())
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := bytes.NewBufferString(tt.requestBody)
			request, err := http.NewRequest(http.MethodPost, testServer.URL+"/value", requestBody)
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Content-Type", "application/json")

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, response.StatusCode)

			var result models.Metrics
			if tt.expectedCode == http.StatusOK {
				// fmt.Printf("response.Body: %v\n", response.Body)
				err := json.NewDecoder(response.Body).Decode(&result)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Printf("result: %v\n", result)

				var buf bytes.Buffer
				encoder := json.NewEncoder(&buf)
				assert.Equal(t, encoder.Encode(result), encoder.Encode(tt.expectedMetrics))
			}

			defer response.Body.Close()
		})
	}
}
