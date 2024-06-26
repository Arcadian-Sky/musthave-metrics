// Пакет handler реализует ручки на все точки коннекта к приложению
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/validate"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/utils"
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
	s   storage.MetricsStorage
	cfg *flags.InitedFlags
}

// NewHandler создает экземпляр Handler
func NewHandler(mStorage storage.MetricsStorage, cnf *flags.InitedFlags) *Handler {
	return &Handler{
		s:   mStorage,
		cfg: cnf,
	}
}

// Получает метрики.
//
// @Summary Получает метрики.
// @Description Обновляет метрику в хранилище.
// @Success 200 {string} string "OK"
// @Router / [get]
func (h *Handler) MetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := validate.CheckHash(validate.GetHashHead(r), body, h.cfg.HashKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Выводим данные
	for name, value := range h.s.GetMetrics(r.Context()) {
		fmt.Fprintf(w, "%s: %v\n", name, value)
	}
}

// Обновляет метрику.
//
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
	body, _ := io.ReadAll(r.Body)
	err := validate.CheckHash(validate.GetHashHead(r), body, h.cfg.HashKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Проверякм переданные параметры
	err = validate.CheckMetricTypeAndName(params.Type, params.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ctx := r.Context()
	// Обновляем метрику
	err = h.s.UpdateMetric(ctx, params.Type, params.Name, params.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Выводим данные
	currentMetrics := h.s.GetMetrics(ctx)
	for name, value := range currentMetrics {
		fmt.Fprintf(w, "%s: %v\n", name, value)
	}
}

// Получает метрику.
//
// @Summary Получает метрику.
// @Description Получает метрику в хранилище.
// @Param type path string true "Тип метрики (gauge или counter)"
// @Param name path string true "Название метрики"
// @Success 200 {string} string "OK"
// @Router /value/{type}/{name} [get]
func (h *Handler) GetMetricHandlerFunc(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := validate.CheckHash(validate.GetHashHead(r), body, h.cfg.HashKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := NewMetricParams(r)
	//Проверякм переданные параметры
	err = validate.CheckMetricTypeAndName(params.Type, params.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//Получаем данные для вывода
	metricTypeID, err := utils.GetMetricTypeByCode(params.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим данные
	fmt.Println("metricTypeID", metricTypeID)
	currentMetrics := h.s.GetMetric(r.Context(), metricTypeID)
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
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		for name, value := range currentMetrics {
			fmt.Fprintf(w, "%s: %v\n", name, value)
		}
	}
}

// Получает метрики через JSON
//
// @Summary Получает метрики.
// @Accept json
// @Produce json
// @Param data body models.Metrics true "Данные в формате JSON"
// @Success 200 {object} string "OK"
// @Failure 404 {object} string "Error"
// @Router /value [post]
func (h *Handler) GetMetricsJSONHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = validate.CheckHash(validate.GetHashHead(r), body, h.cfg.HashKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var metrics models.Metrics

	// Проверяем тело запроса на пустоту
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Декодируем JSON из []byte в структуру Metrics
	if err := json.Unmarshal(body, &metrics); err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		fmt.Printf("Failed to decode JSON:  err.Error(): %v\n", err.Error())
		return
	}

	if metrics.MType != "" && metrics.ID != "" {
		// Выводим данные
		err = h.s.GetJSONMetric(r.Context(), &metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			fmt.Printf("GetJSONMetric err.Error(): %v\n", err.Error())
			return
		}
	}
	resp, err := json.Marshal(&metrics)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка записи Body:", err)
		return
	}
}

// Обновляет метрику через JSON
//
// @Summary Обновляет метрику.
// @Description Обновляет метрику в хранилище через json обьект.
// @Accept json
// @Produce json
// @Param data body models.Metrics true "Данные в формате JSON"
// @Success 200 {object} string "OK"
// @Failure 404 {object} string "Error"
// @Router /update [post]
func (h *Handler) UpdateJSONMetricHandlerFunc(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = validate.CheckHash(validate.GetHashHead(r), body, h.cfg.HashKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var metrics models.Metrics

	// fmt.Printf("body2: %v\n", body)
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
	ctx := r.Context()
	// Обновляем метрику
	err = h.s.UpdateJSONMetric(ctx, &metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим данные
	err = h.s.GetJSONMetric(ctx, &metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(&metrics)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка записи Body:", err)
		return
	}

}

// Пинг БД
func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := validate.CheckHash(validate.GetHashHead(r), body, h.cfg.HashKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.Ping()
	if err != nil {
		http.Error(w, "ошибка при проверке подключения к базе данных:"+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s\n", "Подключение к базе данных успешно!")
	// w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Обновляет метрики через JSON
//
// @Summary Обновляет метрику.
// @Description Обновляет метрику в хранилище через json обьект.
// @Accept json
// @Produce json
// @Param data body models.Metrics true "Данные в формате JSON"
// @Success 200 {object} string "OK"
// @Failure 404 {object} string "Error"
// @Router /updates [post]
func (h *Handler) UpdateJSONMetricsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// err = validate.CheckHash(validate.GetHashHead(r), body, h.cfg.HashKey)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	var metrics []models.Metrics

	// Проверяем тело запроса на пустоту
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// // Декодируем JSON из []byte в структуру Metrics
	if err := json.Unmarshal(body, &metrics); err != nil {
		http.Error(w, "Failed to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Обновляем метрики
	err = h.s.UpdateJSONMetrics(r.Context(), &metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим данные
	// for range metrics{
	// 	err = h.s.GetJSONMetrics(&metrics)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// }

	resp, err := json.Marshal(&metrics)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Ошибка записи Body:", err)
		return
	}
}
