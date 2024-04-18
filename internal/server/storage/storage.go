package storage

import (
	"context"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
)

// MetricType определяет тип метрики (gauge или counter)
type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

// MetricsStorage определяет интерфейс для взаимодействия с хранилищем метрик
type MetricsStorage interface {
	GetMetric(ctx context.Context, mtype MetricType) map[string]interface{}
	UpdateMetric(ctx context.Context, mtype string, name string, value string) error

	GetJSONMetric(ctx context.Context, metric *models.Metrics) error
	UpdateJSONMetric(ctx context.Context, metric *models.Metrics) error

	// GetJSONMetrics(metrics *[]models.Metrics)
	UpdateJSONMetrics(ctx context.Context, metrics *[]models.Metrics) error

	GetMetrics(ctx context.Context) map[MetricType]map[string]interface{}
	SetMetrics(ctx context.Context, metrics map[MetricType]map[string]interface{})

	Ping() error
}
