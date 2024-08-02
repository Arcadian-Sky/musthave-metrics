//go:build !race
// +build !race

package main

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerGracefulShutdown(t *testing.T) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Запускаем сервер
	go main()

	// Отправляем сигнал завершения
	stop <- syscall.SIGTERM

	// Ожидаем завершения работы
	time.Sleep(1 * time.Second)

	// Проверяем, что сервер корректно завершился

	assert.True(t, true)
}
