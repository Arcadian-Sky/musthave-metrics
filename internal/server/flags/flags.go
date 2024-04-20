package flags

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Флаг -a, переменная окружения ADDRESS — endpoint address.
// Флаг -i, переменная окружения STORE_INTERVAL — интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной).
// Флаг -f, переменная окружения FILE_STORAGE_PATH — полное имя файла, куда сохраняются текущие значения (по умолчанию /tmp/metrics-db.json, пустое значение отключает функцию записи на диск).
// Флаг -r, переменная окружения RESTORE — булево значение (true/false), определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера (по умолчанию true).
// Флаг -d, переменная окружения DATABASE_DSN - cтрока с адресом подключения к БД (по умолчанию пусто).

type InitedFlags struct {
	Endpoint       string
	StoreInterval  time.Duration
	FileStorage    string
	RestoreMetrics bool
	DBSettings     string
	StorageType    string
	HashKey        string
}

func Parse() *InitedFlags {
	end := flag.String("a", ":8080", "endpoint address")
	key := flag.String("k", "", "hash key")
	flagStoreInterval := flag.Int("i", 300, "Интервал сохранения метрик на диск")
	flagFileStorage := flag.String("f", "/tmp/metrics-db.json", "Путь к файлу для хранения метрик")
	flagRestoreMetrics := flag.Bool("r", true, "Восстановление метрик при старте сервера")
	flagDBSettings := flag.String("d", "", "Адрес подключения к БД")
	storageType := "inmemory"

	flag.Parse()
	_ = godotenv.Load()

	endpoint := *end
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		endpoint = envRunAddr
	}

	hashKey := *key
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		hashKey = envHashKey
	}

	// fmt.Printf("flagStoreInterval: %v\n", *flagStoreInterval)
	interval := time.Duration(*flagStoreInterval) * time.Second
	if envRunInterv := os.Getenv("STORE_INTERVAL"); envRunInterv != "" {
		if dur, err := time.ParseDuration(envRunInterv); err == nil {
			interval = dur
		}
	}

	fileStorage := *flagFileStorage
	if envRunFileStorage := os.Getenv("FILE_STORAGE_PATH"); envRunFileStorage != "" {
		fileStorage = envRunFileStorage
	}

	restoreMetrics := *flagRestoreMetrics
	if envRunFileStorage := os.Getenv("RESTORE"); envRunFileStorage != "" {
		if val, err := strconv.ParseBool(envRunFileStorage); err == nil {
			restoreMetrics = val
		}
	}

	dbSettings := *flagDBSettings
	if envRunDBSettings := os.Getenv("DATABASE_DSN"); envRunDBSettings != "" {
		dbSettings = envRunDBSettings
	}
	if dbSettings != "" {
		storageType = "postgres"
	}
	// fmt.Printf("flag interval: %v\n", interval)
	// fmt.Printf("flag fileStorage: %v\n", fileStorage)
	// fmt.Printf("flag restoreMetrics: %v\n", restoreMetrics)

	return &InitedFlags{
		Endpoint:       endpoint,
		StoreInterval:  interval,
		FileStorage:    fileStorage,
		RestoreMetrics: restoreMetrics,
		DBSettings:     dbSettings,
		StorageType:    storageType,
		HashKey:        hashKey,
	}

}
