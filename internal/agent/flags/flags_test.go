package flags

import (
	"os"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// Mock environment variables
	os.Setenv("ADDRESS", "testhost:8080")
	os.Setenv("REPORT_INTERVAL", "5")
	os.Setenv("POLL_INTERVAL", "15")

	Parse()

	// Check if the values were correctly parsed
	if serverAddress != "http://testhost:8080" {
		t.Errorf("Expected serverAddress to be 'http://testhost:8080', got '%s'", serverAddress)
	}

	expectedReportInterval := 5 * time.Second
	if reportInterval != expectedReportInterval {
		t.Errorf("Expected reportInterval to be %v, got %v", expectedReportInterval, reportInterval)
	}

	expectedPollInterval := 15 * time.Second
	if pollInterval != expectedPollInterval {
		t.Errorf("Expected pollInterval to be %v, got %v", expectedPollInterval, pollInterval)
	}
}

func TestParseWithDefaults(t *testing.T) {
	// Reset environment variables
	os.Unsetenv("ADDRESS")
	os.Unsetenv("REPORT_INTERVAL")
	os.Unsetenv("POLL_INTERVAL")

	// Parse without setting environment variables
	Parse()

	// Check if the default values were set
	if serverAddress != "http://localhost:8080" {
		t.Errorf("Expected serverAddress to be 'http://localhost:8080', got '%s'", serverAddress)
	}

	expectedReportInterval := 2 * time.Second
	if reportInterval != expectedReportInterval {
		t.Errorf("Expected reportInterval to be %v, got %v", expectedReportInterval, reportInterval)
	}

	expectedPollInterval := 10 * time.Second
	if pollInterval != expectedPollInterval {
		t.Errorf("Expected pollInterval to be %v, got %v", expectedPollInterval, pollInterval)
	}
}
