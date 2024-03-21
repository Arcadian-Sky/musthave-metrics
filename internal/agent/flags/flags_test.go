package flags

import (
	"flag"
	"os"
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
