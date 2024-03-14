package metricrepo

import "errors"

type MetricRepository interface {
	SaveMetric(key string, value interface{}) error
	GetMetric(key string) (interface{}, error)
	GetAllMetrics() map[string]interface{}
	DeleteMetric(key string) error
}

type InMemoryMetricRepository struct {
	metrics map[string]interface{}
}

func NewInMemoryMetricRepository() *InMemoryMetricRepository {
	return &InMemoryMetricRepository{
		metrics: make(map[string]interface{}),
	}
}

func (repo *InMemoryMetricRepository) SaveMetric(key string, value interface{}) error {
	repo.metrics[key] = value
	return nil
}

func (repo *InMemoryMetricRepository) GetMetric(key string) (interface{}, error) {
	metric, ok := repo.metrics[key]
	if !ok {
		return nil, errors.New("metric not found")
	}
	return metric, nil
}

func (repo *InMemoryMetricRepository) GetAllMetrics() map[string]interface{} {
	return repo.metrics
}

func (repo *InMemoryMetricRepository) DeleteMetric(key string) error {
	delete(repo.metrics, key)
	return nil
}
