package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/controller"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	// Создаем канал для обработки сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Загружаем конфигурацию агента
	config, err := flags.Parse()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	agent := controller.NewCollectAndSendMetricsService(&config)

	// Запуск агента в отдельной горутине
	agent.Run()

	<-stop

	// Корректное завершение работы агента
	log.Println("Shutting down agent gracefully...")
	agent.Stop()

	log.Println("Agent stopped gracefully")
}
