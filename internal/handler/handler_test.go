package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricsHandler(t *testing.T) {
	// Создаем фейковое хранилище
	storage.Storage = storage.NewMemStorage()

	// Создаем HTTP запрос для обновления метрики
	req, err := http.NewRequest("POST", "/update/gauge/example_metric/42", nil)
	assert.NoError(t, err)

	// Создаем ResponseRecorder для записи ответа сервера
	rr := httptest.NewRecorder()

	// Создаем обработчик и вызываем его метод ServeHTTP
	UpdateMetricsHandler().ServeHTTP(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	assert.Contains(t, rr.Body.String(), "0: map[example_metric:42]\n")
}

func TestUpdateMetricsHandlerBadRequest(t *testing.T) {
	// Создаем фейковое хранилище
	storage.Storage = storage.NewMemStorage()

	// Создаем HTTP запрос с некорректным типом метрики
	req, err := http.NewRequest("POST", "/update/invalid_type/example_metric/42", nil)
	assert.NoError(t, err)

	// Создаем ResponseRecorder для записи ответа сервера
	rr := httptest.NewRecorder()

	// Создаем обработчик и вызываем его метод ServeHTTP
	UpdateMetricsHandler().ServeHTTP(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Проверяем тело ответа
	// assert.Contains(t, rr.Body.String(), "Metric type validation failed")
}

func TestUpdateMetricsHandlerNotFound(t *testing.T) {
	// Создаем фейковое хранилище
	storage.Storage = storage.NewMemStorage()

	// Создаем HTTP запрос без имени метрики
	req, err := http.NewRequest("POST", "/update/gauge/", nil)
	assert.NoError(t, err)

	// Создаем ResponseRecorder для записи ответа сервера
	rr := httptest.NewRecorder()

	// Создаем обработчик и вызываем его метод ServeHTTP
	UpdateMetricsHandler().ServeHTTP(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Проверяем тело ответа
	assert.Contains(t, rr.Body.String(), "Metric name not provided")
}
