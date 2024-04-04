package config

import (
	"fmt"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/caretaker"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

//	Интервал с которым производим бэкапирование: parsed.StoreInterval,
//	Место в которое производим бэкапирование: parsed.FileStorage,
//	Делаем ли восстановление значений при старте сервера:parsed.RestoreMetrics,
//
// InitConfig инициализирует конфигурацию сервера на основе переданных параметров
func InitConfig(storeMetrics *storage.MemStorage, config *flags.InitedFlags) error {
	// Создаем экземпляр Caretaker для работы с мементо
	caretaker := &caretaker.Caretaker{}
	if config.FileStorage == "" {
		return nil
	}

	// Если задана периодичность, создаем таймер для сохранения данных
	if config.StoreInterval > 0 {
		go func() {
			ticker := time.NewTicker(config.StoreInterval)
			defer ticker.Stop()

			for {
				<-ticker.C
				SaveMetricsToFile(storeMetrics, config.FileStorage)
			}

		}()
	} else {
		// Если периодичность равна 0, выполняем сохранение синхронно
		SaveMetricsToFile(storeMetrics, config.FileStorage)
	}
	// Если задан флаг восстановления, выполняем восстановление значений
	if config.RestoreMetrics {
		restoredMemento, err := caretaker.ReadFromFile(config.FileStorage)
		if err != nil {
			return fmt.Errorf("error reading memento from file: %v", err)
		}
		if restoredMemento != nil {
			storeMetrics.RestoreFromMemento(restoredMemento)
		}
	}

	return nil
}

// saveMetricsToFile сохраняет текущие значения метрик на диск
func SaveMetricsToFile(storeMetrics *storage.MemStorage, fileStoragePath string) {
	memento := storeMetrics.CreateMemento()
	caretaker := &caretaker.Caretaker{}
	err := caretaker.SaveToFile(memento, fileStoragePath)
	if err != nil {
		fmt.Println("Error saving memento:", err)
		// Обработка ошибок сохранения
	}
}
