package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Send(t *testing.T) {
	metrics := make(map[string]interface{})
	pollCount := 10 // Пример значения pollCount

	err := send(metrics, pollCount)
	assert.Error(t, err, "Функция должна вернуть errоr")
}
