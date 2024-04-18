package persistent

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/mock"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/utils"
	"github.com/golang/mock/gomock"
)

func TestGetMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricsStorage := mock.NewMockMetricsStorage(ctrl)
	ctx := context.Background()
	expectedResult := map[string]interface{}{
		"metric1": 100,
		"metric2": 200,
	}
	mtype := "gauge"
	var storeType storage.MetricType

	storeType, err := utils.GetMetricTypeByCode(mtype)
	if err != nil {
		t.Errorf("Expected result no err got %v", err.Error())
	}
	mockMetricsStorage.EXPECT().GetMetric(ctx, storeType).Return(expectedResult)

	result := mockMetricsStorage.GetMetric(ctx, storeType)

	// Добавьте проверку результата
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result %v, got %v", expectedResult, result)
	}
}

func TestGetJSONMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricsStorage := mock.NewMockMetricsStorage(ctrl)
	ctx := context.Background()

	expectedMetric := &models.Metrics{
		ID:    "metric1",
		MType: "gauge",
	}
	var expectedValue float64 = 100
	expectedMetric.Value = &expectedValue

	expectedJSON, err := json.Marshal(expectedMetric)
	if err != nil {
		t.Errorf("Failed to marshal expected metric: %v", err)
	}

	// Задаем ожидаемое возвращаемое значение
	mockMetricsStorage.EXPECT().GetJSONMetric(ctx, gomock.Any()).Return(nil).Do(func(ctx context.Context, metric *models.Metrics) {
		*metric = *expectedMetric
	})

	var modelMetrics = &models.Metrics{}
	// Вызываем метод и проверяем результат
	err = mockMetricsStorage.GetJSONMetric(ctx, modelMetrics)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Преобразуем пустую метрику в JSON
	modelMetricsJSON, err := json.Marshal(modelMetrics)
	if err != nil {
		t.Errorf("Failed to marshal empty metric: %v", err)
	}

	// Сравниваем JSON представления метрик
	if string(expectedJSON) != string(modelMetricsJSON) {
		t.Errorf("Expected metric %s, got %s", expectedJSON, modelMetricsJSON)
	}
}
