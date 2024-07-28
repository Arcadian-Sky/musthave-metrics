package flags

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	serverAddress  string
	cryptoKey      string
	hashKey        string
	configFilePath string
	pollInterval   time.Duration
	reportInterval time.Duration
	rateLimit      int
}
type AgentConfig struct {
	ServerAddress  string `json:"server_address"`
	PollInterval   int    `json:"poll_interval"`
	ReportInterval int    `json:"report_interval"`
	CryptoKey      string `json:"crypto_key"`
}

// Дефолтные значения для теста
func SetDefault() *Config {
	return &Config{
		serverAddress:  "http://localhost:8080",
		reportInterval: time.Second,
		pollInterval:   time.Second,
		rateLimit:      10,
		cryptoKey:      "",
	}
}

// Через флаг -l=<ЗНАЧЕНИЕ> и переменную окружения RATE_LIMIT. - количество одновременно исходящих запросов на сервер нужно ограничивать «сверху»
func Parse() (Config, error) {
	end := flag.String("a", "", "endpoint")
	key := flag.String("k", "", "hash key")
	cryptoKeyPath := flag.String("crypto-key", "", "crypto-key")
	repI := flag.Int("r", 0, "reportInterval")
	polI := flag.Int("p", 0, "pollInterval")
	rLim := flag.Int("l", 0, "rateLimit")
	configFileFlag := flag.String("c", "", "Путь к файлу конфигурации JSON")

	flag.Parse()

	config := Config{
		hashKey:        *key,
		configFilePath: *configFileFlag,
		rateLimit:      *rLim,
	}
	envRunAddr := os.Getenv("ADDRESS")
	cryptoKeyEnv := os.Getenv("CRYPTO_KEY")
	envRepI := os.Getenv("REPORT_INTERVAL")
	envPolI := os.Getenv("POLL_INTERVAL")

	if envConfigFilePath := os.Getenv("CONFIG"); envConfigFilePath != "" {
		config.configFilePath = envConfigFilePath
	}

	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		config.hashKey = envHashKey
	}

	if envRLim := os.Getenv("RATE_LIMIT"); envRLim != "" {
		rateLimit, err := strconv.Atoi(envRLim)
		if err != nil {
			fmt.Println("Error parsing REPORT_INTERVAL:", err)
			return config, err
		}
		config.rateLimit = rateLimit
	}

	var fileConfig AgentConfig
	if config.configFilePath != "" {
		if err := fileConfig.LoadConfig(config.configFilePath); err != nil {
			return config, fmt.Errorf("ошибка загрузки конфигурации из файла: %v", err)
		}
	}

	config.cryptoKey = getString(*cryptoKeyPath, cryptoKeyEnv, fileConfig.CryptoKey, "", "")
	prefix := "http://"
	if config.cryptoKey != "" {
		prefix = "https://"
	}
	config.serverAddress = getString(*end, envRunAddr, fileConfig.ServerAddress, "localhost:8080", prefix)
	config.pollInterval = getDurationFromInt(*polI, envPolI, fileConfig.PollInterval, 10)
	config.reportInterval = getDurationFromInt(*repI, envRepI, fileConfig.ReportInterval, 2)

	return config, nil
}

func getString(flagValue string, envValue string, fileValue string, defaultValue string, prefix string) string {
	if envValue != "" {
		return prefix + envValue
	}
	if flagValue != "" {
		return prefix + flagValue
	}
	if fileValue != "" {
		return prefix + fileValue
	}
	return prefix + defaultValue
}

func getDurationFromInt(flagValue int, envValue string, fileValue int, defaultValue int) time.Duration {
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

func (a *AgentConfig) LoadConfig(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = json.Unmarshal(file, a)
	if err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	return nil
}

func (c *Config) SetConfigServer(s string) {
	c.serverAddress = s
}

func (c *Config) GetServerAddress() string {
	return c.serverAddress
}

func (c *Config) GetCryptoKeyPath() (*rsa.PublicKey, bool) {
	if c.cryptoKey != "" {
		publicKey, err := c.loadCryptoKey(c.cryptoKey)
		if err != nil {
			fmt.Println("Ошибка при загрузке публичного ключа:", err)
			return nil, false
		}
		return publicKey, true
	}

	return nil, false
}

func (c *Config) GetRateLimit() int {
	return c.rateLimit
}

func (c *Config) GetHash() string {
	return c.hashKey
}

func (c *Config) GetReportInterval() time.Duration {
	return c.reportInterval
}

func (c *Config) GetPollInterval() time.Duration {
	return c.pollInterval
}

// loadPublicKey загружает публичный ключ из PEM файла
func (c *Config) loadCryptoKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("неверный формат ключа")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("неверный тип ключа")
	}

	return rsaPublicKey, nil
}
