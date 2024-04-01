package flags

import (
	"flag"
	"os"
	"testing"
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
}
