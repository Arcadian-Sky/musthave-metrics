package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/validate"
)

// Сборщик параметров
type MetricParams struct {
	Type  string
	Name  string
	Value string
}

// NewMetricParams создает экземпляр MetricParams из объекта *http.Request
func NewMetricParams(r *http.Request) MetricParams {
	return MetricParams{
		Type:  chi.URLParam(r, "type"),
		Name:  chi.URLParam(r, "name"),
		Value: chi.URLParam(r, "value"),
	}
}

// Server handlers
type Handler struct {
	s storage.MetricsStorage
}

// NewHandler создает экземпляр Handler
func NewHandler(mStorage storage.MetricsStorage) *Handler {
	return &Handler{
		s: mStorage,
	}
}

// Получает метрики.
// @Summary Получает метрики.
// @Description Обновляет метрику в хранилище.
// @Success 200 {string} string "OK"
// @Router / [get]
func (h *Handler) MetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Выводим данные
	for name, value := range h.s.GetMetrics() {
		fmt.Fprintf(w, "%d: %v\n", name, value)
	}
}

// Обновляет метрику.
// @Summary Обновляет метрику.
// @Description Обновляет метрику в хранилище.
// @Param type path string true "Тип метрики (gauge или counter)"
// @Param name path string true "Название метрики"
// @Param value path string true "Значение метрики"
// @Router /update/{type} [post]
// @Failure 404 {string} string "metric name not provided"
// @Router /update/{type}/{name} [post]
// @Failure 404 {string} string "metric value not provided"
// @Router /update/{type}/{name}/{value} [post]
// @Success 200 {string} string "OK"
// @Failure 404 {string} string "Error"
func (h *Handler) UpdateMetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	params := NewMetricParams(r)

	//Проверякм переданные параметры
	err := validate.CheckMetricTypeAndName(params.Type, params.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Обновляем метрику
	err = h.s.UpdateMetric(params.Type, params.Name, params.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим данные
	currentMetrics := h.s.GetMetrics()
	for name, value := range currentMetrics {
		fmt.Fprintf(w, "%d: %v\n", name, value)
	}

	w.WriteHeader(http.StatusOK)
}

// Получает метрику.
// @Summary Получает метрику.
// @Description Получает метрику в хранилище.
// @Param type path string true "Тип метрики (gauge или counter)"
// @Param name path string true "Название метрики"
// @Success 200 {string} string "OK"
// @Router /value/{type}/{name} [get]
func (h *Handler) GetMetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	params := NewMetricParams(r)

	//Проверякм переданные параметры
	err := validate.CheckMetricTypeAndName(params.Type, params.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//Получаем данные для вывода
	metricTypeID, err := storage.GetMetricTypeByCode(params.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим данные
	currentMetrics := h.s.GetMetric(metricTypeID)
	fmt.Printf("metricName: %v\n", params.Name)
	if params.Name != "" {
		fmt.Printf("currentMetrics[metricName]: %v\n", currentMetrics[params.Name])
		if currentMetrics[params.Name] != nil {
			_, err = w.Write([]byte(fmt.Sprintf("%v", currentMetrics[params.Name])))
			if err != nil {
				http.Error(w, "w.Write Error: "+err.Error(), http.StatusNotFound)
			}
		} else {
			http.Error(w, "Metric value not provided", http.StatusNotFound)
		}
	} else {
		for name, value := range currentMetrics {
			fmt.Fprintf(w, "%s: %v\n", name, value)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Получает метрики.
// @Accept json
// @Produce json
// @Param data body models.Metrics true "Данные в формате JSON"
// @Success 200 {object} string "OK"
// @Failure 404 {object} string "Error"
// @Router /value [post]
func (h *Handler) GetMetricsJSONHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var metrics models.Metrics

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем тело запроса на пустоту
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Декодируем JSON из []byte в структуру Metrics
	if err := json.Unmarshal(body, &metrics); err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим данные
	err = h.s.GetJSONMetric(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(&metrics)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка записи Body:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Обновляет метрику.
// @Description Обновляет метрику в хранилище через json обьект.
// @Accept json
// @Produce json
// @Param data body models.Metrics true "Данные в формате JSON"
// @Success 200 {object} string "OK"
// @Failure 404 {object} string "Error"
// @Router /update [post]
func (h *Handler) UpdateJSONMetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var metrics models.Metrics

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем тело запроса на пустоту
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Декодируем JSON из []byte в структуру Metrics
	if err := json.Unmarshal(body, &metrics); err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("metrics: %v\n", metrics)
	// Обновляем метрику
	err = h.s.UpdateJSONMetric(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим данные
	err = h.s.GetJSONMetric(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("metrics: %v\n", *metrics.Delta)

	resp, err := json.Marshal(&metrics)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка записи Body:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
