package flags

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	serverAddress  string
	hashKey        string
	pollInterval   time.Duration
	reportInterval time.Duration
	rateLimit      int
}

func (c *Config) SetConfigServer(s string) {
	c.serverAddress = s
}

func (c *Config) GetServerAddress() string {
	return c.serverAddress
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

// Дефолтные значения для теста
func SetDefault() *Config {
	return &Config{
		serverAddress:  "http://localhost:8080",
		reportInterval: time.Second,
		pollInterval:   time.Second,
		rateLimit:      10,
	}
}

// Через флаг -l=<ЗНАЧЕНИЕ> и переменную окружения RATE_LIMIT. - количество одновременно исходящих запросов на сервер нужно ограничивать «сверху»
func Parse() (Config, error) {
	end := flag.String("a", "localhost:8080", "endpoint")
	key := flag.String("k", "", "hash key")
	repI := flag.Int("r", 2, "reportInterval")
	polI := flag.Int("p", 10, "pollInterval")
	rLim := flag.Int("l", 0, "rateLimit")
	flag.Parse()

	config := Config{
		serverAddress:  "http://" + *end,
		hashKey:        *key,
		rateLimit:      *rLim,
		reportInterval: time.Duration(*repI) * time.Second,
		pollInterval:   time.Duration(*polI) * time.Second,
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		config.serverAddress = "http://" + envRunAddr
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

	if envRepI := os.Getenv("REPORT_INTERVAL"); envRepI != "" {
		durationRepI, err := strconv.Atoi(envRepI)
		if err != nil {
			fmt.Println("Error parsing REPORT_INTERVAL:", err)
			return config, err
		}
		config.reportInterval = time.Duration(durationRepI) * time.Second
	}

	if envPolI := os.Getenv("POLL_INTERVAL"); envPolI != "" {
		durationPolI, err := strconv.Atoi(envPolI)
		if err != nil {
			fmt.Println("Error parsing POLL_INTERVAL:", err)
			return config, err
		}
		config.pollInterval = time.Duration(durationPolI) * time.Second
	}

	return config, nil
}
