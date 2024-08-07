package main

import (
	"context"
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
	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	// Создаем WaitGroup для синхронизации завершения всех горутин

	// Загружаем конфигурацию агента
	config, err := flags.Parse()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	serviceController := controller.NewCollectAndSendMetricsService(&config)
	go serviceController.Run(ctx)

	<-stop
	fmt.Println("Получен сигнал, останавливаем все горутины...")

	// Отменяем контекст, чтобы уведомить все горутины о завершении
	cancel()

	// Ожидаем завершения всех горутин
	serviceController.Wg.Wait()

	log.Println("Agent stopped gracefully")
}
