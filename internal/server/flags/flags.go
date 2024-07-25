package flags

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
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
	// BPprofEnabled  bool
	DBSettings    string
	StorageType   string
	HashKey       string
	cryptoKeyPath string
}

func Parse() *InitedFlags {
	end := flag.String("a", ":8080", "endpoint address")
	cryptoKeyFlag := flag.String("crypto-key", "", "Путь до файла с публичным ключом для шифрования")
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

	cryptoKey := *cryptoKeyFlag
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cryptoKey = envCryptoKey
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

	// var bPprofEnabled = false
	// if envBPprofEnabled := os.Getenv("BPPROF"); envBPprofEnabled != "" {
	// 	if val, err := strconv.ParseBool(envBPprofEnabled); err == nil {
	// 		bPprofEnabled = val
	// 	}
	// }

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
		cryptoKeyPath:  cryptoKey,
		// BPprofEnabled:  bPprofEnabled,
	}

}
func (i *InitedFlags) GetCryptoKey() *rsa.PrivateKey {
	privateKey, err := i.loadPrivateKey(i.cryptoKeyPath)
	if err != nil {
		log.Fatalf("Ошибка при загрузке приватного ключа: %v", err)
		return nil
	}
	return privateKey
}

func (i *InitedFlags) loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("неверный формат ключа")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
