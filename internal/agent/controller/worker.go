package controller

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/repository"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func getSystemInfo(metricsInfoChan chan<- map[string]interface{}, metricsRepo repository.MetricsRepository) {
	metrics, err := metricsRepo.GetMetrics()
	if err != nil {
		fmt.Println("Error collecting metrics:", err)
		return
	}

	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Ошибка при получении информации о памяти:", err)
		return
	}

	cpuCount, err := cpu.Counts(false)
	if err != nil {
		fmt.Println("Ошибка при получении количества CPU:", err)
		return
	}

	metrics["TotalMemory"] = float64(memoryInfo.Total)
	metrics["FreeMemory"] = float64(memoryInfo.Free)
	metrics["CPUutilization1"] = float64(cpuCount)

	metricsInfoChan <- metrics
}

// worker это наш рабочий, который принимает два канала:
// jobs - канал задач, это входные данные для обработки
// results - канал результатов, это результаты работы воркера
func (c *CollectAndSendMetricsService) worker(id int, jobs <-chan int, results chan<- int, metricsRepo *repository.InMemoryMetricsRepository, pollCount int64) {
	for j := range jobs {
		// для наглядности будем выводить какой рабочий начал работу и какую задачу он выполняет
		log.Println("рабочий", id, "начал выполнение задачи", j)

		// Создаем канал для получения системной информации
		// sysInfoChan := make(chan SystemInfo)
		metricsInfoChan := make(chan map[string]interface{})

		// Запускаем горутину для получения системной информации
		go func() {
			getSystemInfo(metricsInfoChan, metricsRepo)
		}()

		// Ожидаем завершения выполнения горутины с помощью select
		select {
		case metrics := <-metricsInfoChan:
			// err := c.send(metrics, pollCount)
			// if err != nil {
			// 	log.Println("Error sending metrics:", err)
			// }
			// для наглядности выводим, что рабочий завершил задачу
			// отправляем результат в канал результатов
			// err = c.sendPack(metrics, pollCount)
			// if err != nil {
			// 	fmt.Println("Error sending metrics:", err)
			// }
			// atomic.AddInt64(&pollCount, 1)
			c.Push(metrics, pollCount)
			log.Println("рабочий", id, "завершил выполнение задачи", j)
			results <- j + 1

		case <-time.After(time.Second * 5):
			// Если получение системной информации заняло более 5 секунд, выводим сообщение об ошибке
			log.Println("Ошибка: превышено время ожидания получения системной информации")
		}
	}
}

func (c *CollectAndSendMetricsService) Push(metrics map[string]interface{}, pollCount int64) {
	err := c.send(metrics, pollCount)
	if err != nil {
		fmt.Println("Error sending metrics:", err)
	}
	atomic.AddInt64(&pollCount, 1)
	err = c.sendPack(metrics, pollCount)
	if err != nil {
		fmt.Println("Error sending metrics:", err)
	}
}

func (c *CollectAndSendMetricsService) Init(metricsRepo *repository.InMemoryMetricsRepository, pollCount int64) {
	numWorkers := c.config.GetRateLimit()
	numJobs := 5
	fmt.Printf("кол-во задач: %v\n", numJobs)
	fmt.Printf("количество рабочих: %v\n", numWorkers)
	// var pollCount int64
	// atomic.AddInt64(&pollCount, 1)
	// создаем буферизованный канал для принятия задач в воркер
	jobs := make(chan int, numJobs)
	// создаем буферизованный канал для отправки результатов
	results := make(chan int, numJobs)

	// создаем и запускаем 3 воркера, это и есть пул,
	// передаем id, это для наглядности, канал задач и канал результатов
	for w := 1; w <= numWorkers; w++ {
		go c.worker(w, jobs, results, metricsRepo, pollCount)
	}

	// в канал задач отправляем какие-то данные
	// задач у нас 5, а воркера RateLimit, одновременно решается только RateLimit задачи
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	// как вы помните, закрываем канал на стороне отправителя
	close(jobs)

	// забираем из канала результатов результаты
	// можно присваивать переменной, или выводить на экран, но мы не будем
	for a := 1; a <= numJobs; a++ {
		<-results
	}
}
