package flags

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestParse tests the Parse function for various configurations.
func TestParse(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		envVars           map[string]string
		want              *InitedFlags
		configFilePath    string
		configFileContent string
	}{
		{
			name: "NoArguments",
			args: []string{},
			envVars: map[string]string{
				"ADDRESS":           "",
				"STORE_INTERVAL":    "",
				"FILE_STORAGE_PATH": "",
				"RESTORE":           "",
				"DATABASE_DSN":      "",
				"CRYPTO_KEY":        "",
				"KEY":               "",
			},
			configFilePath:    "",
			configFileContent: "",
			want: &InitedFlags{
				Endpoint:       ":8080",
				StoreInterval:  300 * time.Second,
				FileStorage:    "/tmp/metrics-db.json",
				RestoreMetrics: false,
				DBSettings:     "",
				StorageType:    "inmemory",
				HashKey:        "",
				CryptoKeyPath:  "",
				ConfigFilePath: "",
			},
		},
		{
			name: "WithArguments",
			args: []string{"-a", ":9090"},
			envVars: map[string]string{
				"ADDRESS":           "",
				"STORE_INTERVAL":    "",
				"FILE_STORAGE_PATH": "",
				"RESTORE":           "",
				"DATABASE_DSN":      "",
				"CRYPTO_KEY":        "",
				"KEY":               "",
			},
			configFilePath:    "",
			configFileContent: "",
			want: &InitedFlags{
				Endpoint:       ":9090",
				StoreInterval:  300 * time.Second,
				FileStorage:    "/tmp/metrics-db.json",
				RestoreMetrics: false,
				DBSettings:     "",
				StorageType:    "inmemory",
				HashKey:        "",
				CryptoKeyPath:  "",
				ConfigFilePath: "",
			},
		},
		{
			name: "WithEnvironmentVariable",
			args: []string{},
			envVars: map[string]string{
				"ADDRESS": "localhost:7070",
			},
			configFilePath:    "",
			configFileContent: "",
			want: &InitedFlags{
				Endpoint:       "localhost:7070", // Expected to match env variable
				StoreInterval:  300 * time.Second,
				FileStorage:    "/tmp/metrics-db.json",
				RestoreMetrics: false,
				DBSettings:     "",
				StorageType:    "inmemory",
				HashKey:        "",
				CryptoKeyPath:  "",
				ConfigFilePath: "",
			},
		},
		{
			name: "WithArgumentsAndEnvironmentVariable",
			args: []string{"-a", ":9090"},
			envVars: map[string]string{
				"ADDRESS": "localhost:7070",
			},
			configFilePath:    "",
			configFileContent: "",
			want: &InitedFlags{
				Endpoint:       "localhost:7070", // Environment variable should take priority
				StoreInterval:  300 * time.Second,
				FileStorage:    "/tmp/metrics-db.json",
				RestoreMetrics: false,
				DBSettings:     "",
				StorageType:    "inmemory",
				HashKey:        "",
				CryptoKeyPath:  "",
				ConfigFilePath: "",
			},
		},
		{
			name: "WithConfigFile",
			args: []string{},
			envVars: map[string]string{
				"CONFIG": "config.json",
			},
			configFilePath:    "config.json",
			configFileContent: `{"address": ":8081", "store_interval": 600, "store_file": "/tmp/config-db.json", "restore": true, "database_dsn": "postgres://localhost", "crypto_key": "/keys/key.pem"}`,
			want: &InitedFlags{
				Endpoint:       ":8081",
				StoreInterval:  600 * time.Second,
				FileStorage:    "/tmp/config-db.json",
				RestoreMetrics: true,
				DBSettings:     "postgres://localhost",
				StorageType:    "postgres",
				HashKey:        "",
				CryptoKeyPath:  "/keys/key.pem",
				ConfigFilePath: "config.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag definitions
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// Set environment variables for the test
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}
			// Set command-line arguments
			os.Args = append([]string{"cmd"}, tt.args...)

			// Create config file if needed
			if tt.configFilePath != "" && tt.configFileContent != "" {
				err := os.WriteFile(tt.configFilePath, []byte(tt.configFileContent), 0644)
				if err != nil {
					t.Fatalf("failed to write config file: %v", err)
				}
				defer os.Remove(tt.configFilePath) // Clean up config file after test
			}

			// Call Parse function and assert the results
			got := Parse()
			assert.Equal(t, tt.want.Endpoint, got.Endpoint)
			assert.Equal(t, tt.want.StoreInterval, got.StoreInterval)
			assert.Equal(t, tt.want.FileStorage, got.FileStorage)
			assert.Equal(t, tt.want.RestoreMetrics, got.RestoreMetrics)
			assert.Equal(t, tt.want.DBSettings, got.DBSettings)
			assert.Equal(t, tt.want.StorageType, got.StorageType)
			assert.Equal(t, tt.want.HashKey, got.HashKey)
			assert.Equal(t, tt.want.CryptoKeyPath, got.CryptoKeyPath)
			assert.Equal(t, tt.want.ConfigFilePath, got.ConfigFilePath)
		})
	}
}

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
			result := getString(tt.flagValue, tt.envValue, tt.fileValue, "")
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name        string
		flagValue   bool
		envValue    string
		fileValue   bool
		expected    bool
		description string
	}{
		{
			name:        "FlagValue set to true",
			envValue:    "",
			fileValue:   false,
			flagValue:   true,
			expected:    true,
			description: "Should return true from flag",
		},
		{
			name:        "EnvValue set to true",
			envValue:    "true",
			fileValue:   false,
			flagValue:   false,
			expected:    true,
			description: "Should return true from env",
		},
		{
			name:        "FileValue set to true",
			envValue:    "",
			fileValue:   true,
			flagValue:   false,
			expected:    true,
			description: "Should return true from file",
		},
		{
			name:        "AllValuesFalse",
			envValue:    "",
			fileValue:   false,
			expected:    false,
			flagValue:   false,
			description: "Should return false as all are false",
		},
		{
			name:        "EnvValue invalid format",
			envValue:    "notabool",
			fileValue:   true,
			flagValue:   false,
			expected:    true,
			description: "Should return file value if env is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBool(tt.flagValue, tt.envValue, tt.fileValue)
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

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for config files
	dir, err := os.MkdirTemp("", "config")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	tests := []struct {
		name          string
		configData    string
		expected      *fileFlags
		expectedError bool
		description   string
	}{
		{
			name: "ValidConfig",
			configData: `{
				"address": "localhost:8080",
				"store_interval": 2,
				"store_file": "/path/to/file.db",
				"restore": true,
				"database_dsn": "user:password@tcp(localhost:3306)/dbname",
				"crypto_key": "/path/to/key.pem"
			}`,
			expected: &fileFlags{
				Endpoint:       "localhost:8080",
				StoreInterval:  2,
				FileStorage:    "/path/to/file.db",
				RestoreMetrics: true,
				DBSettings:     "user:password@tcp(localhost:3306)/dbname",
				CryptoKeyPath:  "/path/to/key.pem",
			},
			expectedError: false,
			description:   "Should load valid config successfully",
		},
		{
			name:          "InvalidConfig",
			configData:    `{"address": "localhost:8080", "restore": "yes"}`,
			expected:      &fileFlags{},
			expectedError: true,
			description:   "Should fail on invalid boolean",
		},
		{
			name:          "EmptyConfig",
			configData:    `{}`,
			expected:      &fileFlags{},
			expectedError: false,
			description:   "Should load empty config",
		},
		{
			name:          "ConfigWithInvalidDuration",
			configData:    `{"store_interval": "invalid_duration"}`,
			expected:      &fileFlags{},
			expectedError: true,
			description:   "Should fail on invalid duration",
		},
		{
			name: "ConfigWithDefaultValues",
			configData: `{
				"store_interval": 0,
				"restore": false
			}`,
			expected: &fileFlags{
				StoreInterval:  0,
				RestoreMetrics: false,
			},
			expectedError: false,
			description:   "Should load config with default values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := fmt.Sprintf("%s/config.json", dir)
			err := os.WriteFile(configPath, []byte(tt.configData), 0644)
			assert.NoError(t, err)

			var config fileFlags
			err = config.LoadConfig(configPath)

			if tt.expectedError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
				assert.Equal(t, tt.expected, &config, tt.description)
			}
		})
	}
}

func TestGetCryptoKey(t *testing.T) {
	// Create a temporary directory for keys
	dir, err := os.MkdirTemp("", "keys")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// Generate a sample private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	// Encode the key to PEM format
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}

	// Write the key to a file
	keyPath := fmt.Sprintf("%s/private.pem", dir)
	keyFile, err := os.Create(keyPath)
	assert.NoError(t, err)
	defer keyFile.Close()

	err = pem.Encode(keyFile, keyPem)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		cryptoKeyPath string
		expectError   bool
		expectNil     bool
		description   string
	}{
		{
			name:          "ValidPrivateKey",
			cryptoKeyPath: keyPath,
			expectNil:     false,
			description:   "Should load a valid private key",
		},
		{
			name:          "InvalidKeyPath",
			cryptoKeyPath: "/invalid/path/key.pem",
			expectError:   true,
			expectNil:     true,
			description:   "Should return nil for invalid path",
		},
		{
			name:          "EmptyKeyPath",
			cryptoKeyPath: "",
			expectNil:     true,
			description:   "Should return nil for empty path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := InitedFlags{CryptoKeyPath: tt.cryptoKeyPath}
			key, _ := config.GetCryptoKey()
			if tt.expectNil {
				assert.Nil(t, key, tt.description)
			} else {
				assert.NotNil(t, key, tt.description)
			}
		})
	}
}

func TestLoadPrivateKey(t *testing.T) {
	// Create a temporary directory for keys
	dir, err := os.MkdirTemp("", "keys")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// Generate a sample private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	// Encode the key to PEM format
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}

	// Write the key to a file
	keyPath := fmt.Sprintf("%s/private.pem", dir)
	keyFile, err := os.Create(keyPath)
	assert.NoError(t, err)
	defer keyFile.Close()

	err = pem.Encode(keyFile, keyPem)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		keyPath     string
		expectNil   bool
		expectError bool
		description string
	}{
		{
			name:        "ValidPrivateKey",
			keyPath:     keyPath,
			expectNil:   false,
			expectError: false,
			description: "Should load a valid private key",
		},
		{
			name:        "InvalidKeyPath",
			keyPath:     "/invalid/path/key.pem",
			expectNil:   true,
			expectError: true,
			description: "Should return error for invalid path",
		},
		{
			name:        "InvalidKeyFormat",
			keyPath:     "",
			expectNil:   true,
			expectError: true,
			description: "Should return error for invalid key format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := InitedFlags{}
			key, err := config.loadPrivateKey(tt.keyPath)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, key, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, key, tt.description)
			}
		})
	}
}
