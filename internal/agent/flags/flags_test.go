package flags

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name        string
		flagValue   string
		envValue    string
		fileValue   string
		expected    string
		description string
	}{
		{
			name:        "FlagValue set",
			envValue:    "",
			flagValue:   "flag",
			fileValue:   "file",
			expected:    "flag",
			description: "Should return flag value",
		},
		{
			name:        "EnvValue set",
			envValue:    "env",
			fileValue:   "file",
			flagValue:   "flag",
			expected:    "env",
			description: "Should return env value",
		},
		{
			name:        "FileValue set",
			envValue:    "",
			flagValue:   "",
			fileValue:   "file",
			expected:    "file",
			description: "Should return file value",
		},
		{
			name:        "AllValuesEmpty",
			flagValue:   "",
			envValue:    "",
			fileValue:   "",
			expected:    "",
			description: "Should return empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getString(tt.flagValue, tt.envValue, tt.fileValue, "", "")
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestGetDurationFromInt(t *testing.T) {
	tests := []struct {
		name        string
		flagValue   int
		envValue    string
		fileValue   int
		expected    time.Duration
		description string
	}{
		{
			name:        "FlagValue set",
			envValue:    "",
			flagValue:   60,
			fileValue:   0,
			expected:    60 * time.Second,
			description: "Should return duration from flag",
		},
		{
			name:        "EnvValue set",
			envValue:    "120",
			flagValue:   0,
			fileValue:   300,
			expected:    120 * time.Second,
			description: "Should return duration from env",
		},
		{
			name:        "FileValue set",
			envValue:    "",
			flagValue:   0,
			fileValue:   300,
			expected:    300 * time.Second,
			description: "Should return duration from file",
		},
		{
			name:        "AllValuesZero",
			flagValue:   0,
			envValue:    "",
			fileValue:   0,
			expected:    0 * time.Second,
			description: "Should return zero duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDurationFromInt(tt.flagValue, tt.envValue, tt.fileValue, 0)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestGetCryptoKeyPathWithFatalLog тестирует метод GetCryptoKeyPath с учетом вызова log.Fatal

// TestGetServerAddress тестирует метод GetServerAddress
func TestGetServerAddress(t *testing.T) {
	config := &Config{serverAddress: "http://localhost:8080"}
	assert.Equal(t, "http://localhost:8080", config.GetServerAddress())

	config.SetConfigServer("http://example.com")
	assert.Equal(t, "http://example.com", config.GetServerAddress())
}

// TestGetRateLimit тестирует метод GetRateLimit
func TestGetRateLimit(t *testing.T) {
	config := &Config{rateLimit: 10}
	assert.Equal(t, 10, config.GetRateLimit())
}

// TestGetHash тестирует метод GetHash
func TestGetHash(t *testing.T) {
	config := &Config{hashKey: "secret"}
	assert.Equal(t, "secret", config.GetHash())
}

// TestGetReportInterval тестирует метод GetReportInterval
func TestGetReportInterval(t *testing.T) {
	config := &Config{reportInterval: 5 * time.Second}
	assert.Equal(t, 5*time.Second, config.GetReportInterval())
}

// TestGetPollInterval тестирует метод GetPollInterval
func TestGetPollInterval(t *testing.T) {
	config := &Config{pollInterval: 10 * time.Second}
	assert.Equal(t, 10*time.Second, config.GetPollInterval())
}

// TestLoadConfig тестирует метод LoadConfig
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name         string
		configData   string
		expected     AgentConfig
		expectErr    bool
		errorMessage string
	}{
		{
			name: "ValidConfig",
			configData: `
				{
					"server_address": "http://localhost:8080",
					"poll_interval": 15,
					"report_interval": 30,
					"crypto_key": "/path/to/crypto.key"
				}
			`,
			expected: AgentConfig{
				ServerAddress:  "http://localhost:8080",
				PollInterval:   15,
				ReportInterval: 30,
				CryptoKey:      "/path/to/crypto.key",
			},
			expectErr: false,
		},
		{
			name: "InvalidJSON",
			configData: `
				{
					"server_address": "http://localhost:8080",
					"poll_interval": 15,
					"report_interval": 30,
					"crypto_key": "/path/to/crypto.key"
			`,
			expected:     AgentConfig{},
			expectErr:    true,
			errorMessage: "error parsing config file:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем временный файл конфигурации
			configFile, err := os.CreateTemp("", "agent_config_*.json")
			assert.NoError(t, err)
			defer os.Remove(configFile.Name())

			// Записываем данные конфигурации в файл
			_, err = configFile.Write([]byte(tt.configData))
			assert.NoError(t, err)

			// Закрываем файл
			err = configFile.Close()
			assert.NoError(t, err)

			// Загружаем конфигурацию
			var agentConfig AgentConfig
			err = agentConfig.LoadConfig(configFile.Name())

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, agentConfig)
			}
		})
	}
}

// TestParseConfig тестирует функцию Parse
func TestParseConfig(t *testing.T) {
	// Устанавливаем флаги и переменные окружения
	os.Setenv("ADDRESS", "http://env_address:8080")
	os.Setenv("CRYPTO_KEY", "/path/to/env_crypto.key")
	os.Setenv("REPORT_INTERVAL", "25")
	os.Setenv("POLL_INTERVAL", "35")
	os.Setenv("RATE_LIMIT", "5")
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("CRYPTO_KEY")
		os.Unsetenv("REPORT_INTERVAL")
		os.Unsetenv("POLL_INTERVAL")
		os.Unsetenv("RATE_LIMIT")
	}()

	// Создаем временный файл конфигурации
	configData := `
	{
		"server_address": "http://file_address:8080",
		"poll_interval": 15,
		"report_interval": 30,
		"crypto_key": "/path/to/file_crypto.key"
	}
	`
	configFile, err := os.CreateTemp("", "agent_config_*.json")
	assert.NoError(t, err)
	defer os.Remove(configFile.Name())

	// Записываем данные конфигурации в файл
	_, err = configFile.Write([]byte(configData))
	assert.NoError(t, err)

	// Закрываем файл
	err = configFile.Close()
	assert.NoError(t, err)

	// Устанавливаем флаги
	os.Args = []string{"cmd", "-a", "http://flag_address:8080", "-k", "flag_hash_key", "-crypto-key", "/path/to/flag_crypto.key", "-r", "10", "-p", "20", "-l", "7", "-c", configFile.Name()}

	// Выполняем разбор конфигурации
	config, err := Parse()
	assert.NoError(t, err)

	// Проверяем значения конфигурации
	assert.Equal(t, "https://http://env_address:8080", config.serverAddress) // Переменная окружения имеет высший приоритет
	assert.Equal(t, "/path/to/env_crypto.key", config.cryptoKey)             // Переменная окружения имеет высший приоритет
	assert.Equal(t, "flag_hash_key", config.hashKey)                         // Флаг имеет высший приоритет
	assert.Equal(t, 5, config.rateLimit)                                     // Переменная окружения имеет высший приоритет
	assert.Equal(t, 25*time.Second, config.reportInterval)                   // Переменная окружения имеет высший приоритет
	assert.Equal(t, 35*time.Second, config.pollInterval)                     // Переменная окружения имеет высший приоритет
}
