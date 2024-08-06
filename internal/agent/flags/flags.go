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
	tEndpoint      string
	tcpEnable      bool
	pollInterval   time.Duration
	reportInterval time.Duration
	rateLimit      int
}

type AgentConfig struct {
	ServerAddress  string       `json:"server_address"`
	PollInterval   JSONDuration `json:"poll_interval"`
	ReportInterval JSONDuration `json:"report_interval"`
	CryptoKey      string       `json:"crypto_key"`
}

type JSONDuration time.Duration

func (d JSONDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *JSONDuration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// Convert the string into a time.Duration
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = JSONDuration(duration)
	return nil
}

// String returns the duration as a string
func (d JSONDuration) String() string {
	return time.Duration(d).String()
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
	tcpEndpoint := flag.String("tcpa", "", "tcp endpoint address")
	tcpEnable := flag.Bool("tb", false, "использовать tcp")

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
	envTCPRunAddr := os.Getenv("TCP_ADDRESS")
	envTCPEnable := os.Getenv("TCP_ENABLE")

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
	config.pollInterval = getDuration(*polI, envPolI, fileConfig.PollInterval, 10)
	config.reportInterval = getDuration(*repI, envRepI, fileConfig.ReportInterval, 2)
	config.tEndpoint = getString(*tcpEndpoint, envTCPRunAddr, "", ":3200", "")
	config.tcpEnable = getBool(*tcpEnable, envTCPEnable, false, false)

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

func getDuration(flagValue int, envValue string, fileValue JSONDuration, defaultValue int) time.Duration {
	if envValue != "" {
		if parsed, err := strconv.Atoi(envValue); err == nil {
			return time.Duration(parsed) * time.Second
		}
	}

	if flagValue != 0 {
		return time.Duration(flagValue) * time.Second
	}

	if fileValue != JSONDuration(0) {
		return time.Duration(fileValue)
	}

	return time.Duration(defaultValue) * time.Second
}

func getBool(flagValue bool, envValue string, fileValue bool, defaultValue bool) bool {
	if envValue != "" {
		if parsed, err := strconv.ParseBool(envValue); err == nil {
			return parsed
		}
	}
	if flagValue {
		return flagValue
	}
	if fileValue {
		return fileValue
	}
	return defaultValue
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

func (c *Config) GetTEndpoint() string {
	return c.tEndpoint
}

func (c *Config) GetTcpEnable() bool {
	return c.tcpEnable
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
