package flags

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
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
	Endpoint       string        `json:"address"`
	StoreInterval  time.Duration `json:"store_interval"`
	FileStorage    string        `json:"store_file"`
	RestoreMetrics bool          `json:"restore"`
	DBSettings     string        `json:"database_dsn"`
	CryptoKeyPath  string        `json:"crypto_key"`
	StorageType    string
	HashKey        string
	ConfigFilePath string
}

func Parse() *InitedFlags {
	address := flag.String("a", "", "endpoint address")
	flagDBSettings := flag.String("d", "", "Адрес подключения к БД")
	flagStoreInterval := flag.Int("i", 0, "Интервал сохранения метрик на диск")
	flagFileStorage := flag.String("f", "", "Путь к файлу для хранения метрик")
	flagRestoreMetrics := flag.Bool("r", false, "Восстановление метрик при старте сервера")
	flagHashKey := flag.String("k", "", "hash key")
	cryptoKeyFlag := flag.String("crypto-key", "", "Путь до файла с публичным ключом для шифрования")
	configFileFlag := flag.String("c", "", "Путь к файлу конфигурации JSON")

	flag.Parse()
	_ = godotenv.Load()

	var initedConfig InitedFlags
	var fileConfig InitedFlags

	envRunAddr := os.Getenv("ADDRESS")
	envRunRestoreStorage := os.Getenv("RESTORE")
	envRunInterv := os.Getenv("STORE_INTERVAL")
	envRunFileStorage := os.Getenv("FILE_STORAGE_PATH")
	envRunDBSettings := os.Getenv("DATABASE_DSN")
	envCryptoKey := os.Getenv("CRYPTO_KEY")
	envHashKey := os.Getenv("KEY")
	configFilePathEnv := os.Getenv("CONFIG")

	configFilePath := *configFileFlag
	if configFilePathEnv != "" {
		configFilePath = configFilePathEnv
	}
	if configFilePath != "" {
		if err := fileConfig.LoadConfig(configFilePath); err != nil {
			log.Fatalf("Ошибка загрузки конфигурации из файла: %v", err)
		}
	}

	fmt.Printf("fileConfig: %v\n", fileConfig)

	initedConfig.ConfigFilePath = getString(*configFileFlag, configFilePathEnv, "", "")
	initedConfig.DBSettings = getString(*flagDBSettings, envRunDBSettings, fileConfig.DBSettings, "")
	initedConfig.Endpoint = getString(*address, envRunAddr, fileConfig.Endpoint, ":8080")
	initedConfig.StoreInterval = getDurationFromInt(*flagStoreInterval, envRunInterv, fileConfig.StoreInterval, 300)
	initedConfig.FileStorage = getString(*flagFileStorage, envRunFileStorage, fileConfig.FileStorage, "/tmp/metrics-db.json")
	initedConfig.CryptoKeyPath = getString(*cryptoKeyFlag, envCryptoKey, fileConfig.CryptoKeyPath, "")
	initedConfig.HashKey = getString(*flagHashKey, envHashKey, "", "")
	initedConfig.RestoreMetrics = getBool(*flagRestoreMetrics, envRunRestoreStorage, fileConfig.RestoreMetrics)

	initedConfig.StorageType = "inmemory"
	if initedConfig.DBSettings != "" {
		initedConfig.StorageType = "postgres"
	}

	return &initedConfig
}

func getString(flagValue string, envValue string, fileValue string, defaultValue string) string {
	if envValue != "" {
		return envValue
	}
	if flagValue != "" {
		return flagValue
	}
	if fileValue != "" {
		return fileValue
	}
	return defaultValue
}

func getBool(flagValue bool, envValue string, fileValue bool) bool {
	if envValue != "" {
		if parsed, err := strconv.ParseBool(envValue); err == nil {
			return parsed
		}
	}
	if flagValue {
		return flagValue
	}
	return fileValue
}

func getDurationFromInt(flagValue int, envValue string, fileValue time.Duration, defaultValue int) time.Duration {
	if envValue != "" {
		if parsed, err := strconv.Atoi(envValue); err == nil {
			return time.Duration(parsed) * time.Second
		}
	}

	if flagValue != 0 {
		return time.Duration(flagValue) * time.Second
	}

	if fileValue != 0 {
		return time.Duration(fileValue) * time.Second
	}

	return time.Duration(defaultValue) * time.Second
}

func (c *InitedFlags) LoadConfig(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	err = json.Unmarshal(file, c)
	if err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}
	return nil
}

func (i *InitedFlags) GetCryptoKey() (*rsa.PrivateKey, error) {

	if i.CryptoKeyPath != "" {
		privateKey, err := i.loadPrivateKey(i.CryptoKeyPath)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при загрузке приватного ключа: %v", err)
		}
		return privateKey, nil
	}
	return nil, nil
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
