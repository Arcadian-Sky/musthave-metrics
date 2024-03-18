package server

import (
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler"
)

func TestInitRouter(t *testing.T) {

	handler := handler.NewHandler()
	handler.InitStorage()

	t.Run("TestInitRouter", func(t *testing.T) {
		InitRouter(*handler)
	})

}
