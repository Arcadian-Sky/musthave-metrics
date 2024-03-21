package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
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

func TestUpdateMetricsHandlers(t *testing.T) {
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
			expectedCode:  404,
			expectedValue: "404",
		},
		{
			name:          "not valid request 2",
			requestPath:   "/update/counter/someName/",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  404,
			expectedValue: "404",
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
			expectedCode:  404,
			expectedValue: "metric type not provided",
		},
	}

	testServer := httptest.NewServer(InitRouter())
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// resp, _ := testRequest(t, ts, "POST", tt.requestPath)
			// defer resp.Body.Close()
			// // fmt.Println(resp.StatusCode)
			// assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// respBody, _ := io.ReadAll(resp.Body)
			// assert.Contains(t, tt.expectedValue, string(respBody))

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
			// assert.Equal(t, test.expectedBody, string(respBody))

		})
	}

}

func InitRouter() chi.Router {

	handler := NewHandler(storage.NewMemStorage())

	// handler := NewHandler()
	// handler.InitStorage()

	r := chi.NewRouter()

	r.Head("/", func(rw http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "text/plain")
	})
	// GET http://localhost:8080/value/counter/testSetGet163

	r.Get("/", handler.MetricsHandlerFunc)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", handler.UpdateMetricsHandlerFunc)
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}/", handler.UpdateMetricsHandlerFunc)
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/", handler.GetMetricsHandlerFunc)
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/", handler.GetMetricsHandlerFunc)
			r.Get("/{name}", handler.GetMetricsHandlerFunc)
			r.Get("/{name}/", handler.GetMetricsHandlerFunc)
		})
	})

	return r

	// return server.InitRouter(*handler)
}
