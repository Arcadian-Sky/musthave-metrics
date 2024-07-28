//go:build !race
// +build !race

package main

import (
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

func TestOpenDatabase(t *testing.T) {
	tests := []struct {
		name        string
		dbSettings  string
		expectError bool
	}{
		{
			name:        "Valid connection",
			dbSettings:  "postgres://user:password@localhost/dbname?sslmode=disable",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := OpenDatabase(tt.dbSettings)

			if (err != nil) != tt.expectError {
				t.Errorf("OpenDatabase() error = %v, expectError %v", err, tt.expectError)
			}

			if db != nil {
				db.Close()
			}
		})
	}
}

type MockStorage struct {
	storage.MetricsStorage
}

// TestInitializeConfig tests the InitializeConfig function
func TestInitializeConfig(t *testing.T) {
	tests := []struct {
		name        string
		memStore    storage.MetricsStorage
		expectError bool
	}{
		{
			name:        "Successful initialization",
			memStore:    &MockStorage{},
			expectError: false,
		},
		{
			name:        "Failed initialization",
			memStore:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &flags.InitedFlags{}

			err, _, _ := InitializeConfig(tt.memStore, parsed)

			if (err != nil) != tt.expectError {
				t.Errorf("InitializeConfig() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
