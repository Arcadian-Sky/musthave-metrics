package server

import (
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestInitRouter(t *testing.T) {
	tests := []struct {
		name string
		want chi.Router
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
