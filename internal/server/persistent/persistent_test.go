package persistent

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/mock"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/utils"
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

func TestGetMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricsStorage := mock.NewMockMetricsStorage(ctrl)
	ctx := context.Background()

	expectedResult := map[storage.MetricType]map[string]interface{}{
		"gauge": {
			"metric1": 100,
			"metric2": 200,
		},
		"counter": {
			"metric1": 100,
			"metric2": 200,
		},
	}

	mockMetricsStorage.EXPECT().GetMetrics(ctx).Return(expectedResult)

	result := mockMetricsStorage.GetMetrics(context.Background())

	// Добавьте проверку результата
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result %v, got %v", expectedResult, result)
	}
}

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricsStorage := mock.NewMockMetricsStorage(ctrl)

	mockMetricsStorage.EXPECT().Ping().Return(nil)

	result := mockMetricsStorage.Ping()

	// Добавьте проверку результата
	if !reflect.DeepEqual(result, nil) {
		t.Errorf("Expected result %v, got %v", nil, result)
	}
}

func TestSetMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricsStorage := mock.NewMockMetricsStorage(ctrl)
	ctx := context.Background()

	metrics := map[storage.MetricType]map[string]interface{}{
		"gauge": {
			"metric1": 100,
			"metric2": 200,
		},
		"counter": {
			"metric1": 100,
			"metric2": 200,
		},
	}

	expectedResult := map[storage.MetricType]map[string]interface{}{
		"gauge": {
			"metric1": 100,
			"metric2": 200,
		},
		"counter": {
			"metric1": 100,
			"metric2": 200,
		},
	}

	mockMetricsStorage.EXPECT().SetMetrics(ctx, expectedResult)

	mockMetricsStorage.SetMetrics(context.Background(), metrics)

	ctrl.Finish()
}

func TestUpdateJSONMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetricsStorage := mock.NewMockMetricsStorage(ctrl)
	ctx := context.Background()
	var int123 = int64(123)
	metrics := models.Metrics{
		ID:    "23123",
		MType: "gauge",
		Delta: &int123,
	}

	expectedResult := models.Metrics{
		ID:    "23123",
		MType: "gauge",
		Delta: &int123,
	}

	mockMetricsStorage.EXPECT().UpdateJSONMetric(ctx, &expectedResult)

	err := mockMetricsStorage.UpdateJSONMetric(context.Background(), &metrics)

	assert.NoError(t, err, "UpdateJSONMetric() should not return an error")

	ctrl.Finish()
}
