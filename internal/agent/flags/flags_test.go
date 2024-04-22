package flags

import (
	"flag"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// Mock environment variables
	os.Setenv("ADDRESS", "testhost:8080")
	os.Setenv("REPORT_INTERVAL", "5")
	os.Setenv("POLL_INTERVAL", "15")

	config, _ := Parse()

	// Check if the values were correctly parsed
	if config.serverAddress != "http://testhost:8080" {
		t.Errorf("Expected serverAddress to be 'http://testhost:8080', got '%s'", config.serverAddress)
	}

	expectedReportInterval := 5 * time.Second
	if config.reportInterval != expectedReportInterval {
		t.Errorf("Expected reportInterval to be %v, got %v", expectedReportInterval, config.reportInterval)
	}

	expectedPollInterval := 15 * time.Second
	if config.pollInterval != expectedPollInterval {
		t.Errorf("Expected pollInterval to be %v, got %v", expectedPollInterval, config.pollInterval)
	}

	// Очищаем флаги перед каждым тестом
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestSetDefault(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Test with default values",
			want: &Config{
				serverAddress:  "http://localhost:8080",
				reportInterval: time.Second,
				pollInterval:   time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetConfigServer(t *testing.T) {
	config := &Config{}
	expectedServerAddress := "http://example.com"

	config.SetConfigServer(expectedServerAddress)

	if config.serverAddress != expectedServerAddress {
		t.Errorf("Expected server address %s, got %s", expectedServerAddress, config.serverAddress)
	}
}

func TestConfig_GetPollInterval(t *testing.T) {
	tests := []struct {
		name   string
		fields Config
		want   time.Duration
	}{
		{
			name:   "Test with default poll interval",
			fields: *SetDefault(),
			want:   time.Second,
		},
		{
			name: "Test with custom poll interval",
			fields: Config{
				pollInterval: 20 * time.Second,
			},
			want: 20 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				serverAddress:  tt.fields.serverAddress,
				pollInterval:   tt.fields.pollInterval,
				reportInterval: tt.fields.reportInterval,
			}
			if got := c.GetPollInterval(); got != tt.want {
				t.Errorf("Config.GetPollInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetReportInterval(t *testing.T) {
	tests := []struct {
		name   string
		fields Config
		want   time.Duration
	}{
		{
			name:   "Test with default report interval",
			fields: *SetDefault(),
			want:   time.Second,
		},
		{
			name: "Test with custom report interval",
			fields: Config{
				reportInterval: 20 * time.Second,
			},
			want: 20 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				serverAddress:  tt.fields.serverAddress,
				pollInterval:   tt.fields.pollInterval,
				reportInterval: tt.fields.reportInterval,
			}
			if got := c.GetReportInterval(); got != tt.want {
				t.Errorf("Config.GetReportInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetServerAddress(t *testing.T) {
	tests := []struct {
		name   string
		fields Config
		want   string
	}{
		{
			name:   "Test with default server address",
			fields: *SetDefault(),
			want:   "http://localhost:8080",
		},
		{
			name: "Test with custom server address",
			fields: Config{
				serverAddress: "http://example.com",
			},
			want: "http://example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				serverAddress:  tt.fields.serverAddress,
				pollInterval:   tt.fields.pollInterval,
				reportInterval: tt.fields.reportInterval,
			}
			if got := c.GetServerAddress(); got != tt.want {
				t.Errorf("Config.GetServerAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
