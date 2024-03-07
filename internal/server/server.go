package server

import (
	"log"
	"net/http"

	"github.com/Arcadian-Sky/musthave-metrics/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Middleware func(http.Handler) http.Handler

// сontentTypeCheckerMiddleware возвращает middleware, которое проверяет тип данных и устанавливает тип для ответа
func сontentTypeCheckerMiddleware(expectedContentType string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Устанавливаем Content-Type для ответа
			w.Header().Set("Content-Type", expectedContentType)

			// Вызываем следующий обработчик в цепочке
			next.ServeHTTP(w, r)
		})
	}
}

func InitServer() {

	r := InitRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}

func InitRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(сontentTypeCheckerMiddleware("text/plain"))

	r.Head("/", func(rw http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "text/plain")
	})

	r.Get("/", handler.MetricsHandlerFunc)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", handler.UpdateMetricsHandlerFunc)
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}/", handler.UpdateMetricsHandlerFunc)
		})
	})
	return r
}
