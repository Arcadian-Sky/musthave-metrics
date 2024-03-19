package repository

import (
	"testing"
)

func TestInMemoryMetricsRepository_GetMetrics(t *testing.T) {

	t.Run("GetMetrics", func(t *testing.T) {
		metricsRepo := NewInMemoryMetricsRepository()
		got, err := metricsRepo.GetMetrics()

		if (err != nil) != false {
			t.Errorf("InMemoryMetricsRepository.GetMetrics() error = %v, wantErr %v", err, false)
			return
		}

		// fmt.Println(got)
		_, ok := got["BuckHashSys"]
		if !ok {
			t.Errorf("InMemoryMetricsRepository.GetMetrics() key 'BuckHashSys' not found in metrics")
			return
		}

	})
}
