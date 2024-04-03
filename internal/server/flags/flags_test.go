package flags

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		args []string
		env  string
		want string
	}{
		{
			name: "NoArguments",
			args: []string{},
			env:  "",
			want: ":8080",
		},
		{
			name: "WithArguments",
			args: []string{"-a", ":9090"},
			env:  "",
			want: ":9090",
		},
		{
			name: "WithEnvironmentVariable",
			args: []string{},
			env:  "localhost:7070",
			want: "localhost:7070",
		},
		{
			name: "WithArgumentsAndEnvironmentVariable",
			args: []string{"-a", ":9090"},
			env:  "localhost:7070",
			want: "localhost:7070",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// Устанавливаем аргументы командной строки
			os.Args = append([]string{"test"}, tt.args...)

			// Устанавливаем переменную окружения
			os.Setenv("ADDRESS", tt.env)

			// if got := Parse(); got != tt.want {
			// 	t.Errorf("Parse() = %v, want %v", got, tt.want)
			// }
		})
	}

	testsArgs := []struct {
		name               string
		args               []string
		envInterv          string
		envFileStorage     string
		envRestore         string
		flagStoreInterval  int
		flagFileStorage    string
		flagRestoreMetrics bool
		wantStoreInterval  time.Duration
		wantFileStorage    string
		wantRestoreMetrics bool
	}{
		{
			name:               "NoArguments",
			args:               []string{},
			envInterv:          "",
			envFileStorage:     "",
			envRestore:         "",
			flagStoreInterval:  300,
			flagFileStorage:    "/tmp/metrics-db.json",
			flagRestoreMetrics: true,
			wantStoreInterval:  300 * time.Second,
			wantFileStorage:    "/tmp/metrics-db.json",
			wantRestoreMetrics: true,
		},
		{
			name:               "WithEnvironmentVariables",
			args:               []string{},
			envInterv:          "5m",
			envFileStorage:     "/custom/path",
			envRestore:         "false",
			flagStoreInterval:  300,
			flagFileStorage:    "/tmp/metrics-db.json",
			flagRestoreMetrics: true,
			wantStoreInterval:  5 * time.Minute,
			wantFileStorage:    "/custom/path",
			wantRestoreMetrics: false,
		},
	}

	for _, tt := range testsArgs {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = append([]string{"test"}, tt.args...)

			os.Setenv("STORE_INTERVAL", tt.envInterv)
			os.Setenv("FILE_STORAGE_PATH", tt.envFileStorage)
			os.Setenv("RESTORE", tt.envRestore)

			got := Parse()

			if got.StoreInterval != tt.wantStoreInterval {
				t.Errorf("Parse().StoreInterval = %v, want %v", got.StoreInterval, tt.wantStoreInterval)
			}

			if got.FileStorage != tt.wantFileStorage {
				t.Errorf("Parse().FileStorage = %v, want %v", got.FileStorage, tt.wantFileStorage)
			}

			if got.RestoreMetrics != tt.wantRestoreMetrics {
				t.Errorf("Parse().RestoreMetrics = %v, want %v", got.RestoreMetrics, tt.wantRestoreMetrics)
			}
		})
	}
}
