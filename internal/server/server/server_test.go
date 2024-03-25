package server

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

func TestInitRouter(t *testing.T) {
	fakeHandler := handler.NewHandler(storage.NewMemStorage())
	// Получаем роутер с помощью InitRouter
	router := InitRouter(*fakeHandler)
	expectedPaths := []string{
		"/",
		"/update/",
		"/update/{type}/",
		"/update/{type}/{name}/",
		"/update/{type}/{name}/{value}/",
		"/value/",
		"/value/{type}/",
		"/value/{type}/{name}/",
	}
	foundPaths := make(map[string]bool)

	// Используем Walk для обхода всех маршрутов и проверки, что они содержатся в ожидаемом списке
	err := chi.Walk(router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Добавляем путь в мапу найденных путей
		foundPaths[route] = true
		return nil
	})

	if err != nil {
		t.Errorf("Error in chi.Walk %q", err.Error())
	}

	// Проверяем, что все ожидаемые пути были найдены
	for _, path := range expectedPaths {
		if !foundPaths[path] {
			t.Errorf("Expected path %q was not found", path)
		}
	}

	// Проверяем, что нет ненайденных путей
	// for path := range foundPaths {
	// 	if !contains(expectedPaths, path) {
	// 		t.Errorf("Unexpected path %q was found", path)
	// 	}
	// }

}

// func contains(s []string, e string) bool {
// 	for _, a := range s {
// 		if a == e {
// 			return true
// 		}
// 	}
// 	return false
// }
