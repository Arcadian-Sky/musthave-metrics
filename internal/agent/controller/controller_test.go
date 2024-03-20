package controller

import (
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	"github.com/stretchr/testify/assert"
)

func Test_Send(t *testing.T) {
	metrics := make(map[string]interface{})
	pollCount := 10 // Пример значения pollCount

	service := NewCollectAndSendMetricsService(*flags.SetDefault())
	err := service.send(metrics, pollCount)
	assert.Error(t, err, "Функция должна вернуть errоr")
}
