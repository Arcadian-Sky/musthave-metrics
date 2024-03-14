package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	ts := httptest.NewServer(InitRouter())
	defer ts.Close()

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
			expectedValue: "0: map[someName:100.001]",
		},
		{
			name:          "valid request count",
			requestPath:   "/update/counter/someName/100",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  200,
			expectedValue: "0: map[someName:100]",
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
			name:          "not valid request",
			requestPath:   "/update/gauge/someName/",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  404,
			expectedValue: "Metric value not provided",
		},
		{
			name:          "not valid request",
			requestPath:   "/update/counter/someName/",
			expectedType:  "count",
			expectedName:  "someName",
			expectedCode:  404,
			expectedValue: "Metric value not provided",
		},
		{
			name:          "not valid request count no name",
			requestPath:   "/update/counter/",
			expectedType:  "count",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "Metric name not provided",
		},
		{
			name:          "not valid request gauge no name",
			requestPath:   "/update/gauge/",
			expectedType:  "gauge",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "Metric name not provided",
		},
		{
			name:          "not valid request gauge no type",
			requestPath:   "/update/",
			expectedType:  "",
			expectedName:  "",
			expectedCode:  404,
			expectedValue: "Metric type not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, _ := testRequest(t, ts, "POST", tt.requestPath)
			defer resp.Body.Close()
			// fmt.Println(resp.StatusCode)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// actualType := chi.URLParam(req, "type")
			// assert.Equal(t, expectedType, actualType)
			// fmt.Printf("actualType: %v\n", actualType)

			respBody, _ := io.ReadAll(resp.Body)
			assert.Contains(t, tt.expectedValue, string(respBody))

		})
	}

}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, *http.Request) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, req //string(respBody)
}

func InitRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Head("/", func(rw http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "text/plain")
	})

	r.Get("/", MetricsHandlerFunc)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", UpdateMetricsHandlerFunc)
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", UpdateMetricsHandlerFunc)
			r.Post("/{name}", UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}", UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}/", UpdateMetricsHandlerFunc)
		})
	})

	return r
}
