package server

import (
	"reflect"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler"
	"github.com/go-chi/chi/v5"
)

func Test_InitRouter(t *testing.T) {
	type args struct {
		handler handler.Handler
	}
	tests := []struct {
		name string
		args args
		want chi.Router
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitRouter(tt.args.handler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_сontentTypeCheckerMiddleware(t *testing.T) {
	type args struct {
		expectedContentType string
	}
	tests := []struct {
		name string
		args args
		want Middleware
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := сontentTypeCheckerMiddleware(tt.args.expectedContentType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("сontentTypeCheckerMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
