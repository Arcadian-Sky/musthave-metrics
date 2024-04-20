package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Arcadian-Sky/musthave-metrics/docs"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler"
	packmiddleware "github.com/Arcadian-Sky/musthave-metrics/internal/server/middleware"
)

// @title           API
// @version         1.0
// @openapi         3.1
// @description     This is a sample server celler server.
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func InitRouter(handler handler.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(packmiddleware.Logger)

	// r.Use(middleware.Logger)
	// r.Use(packmiddleware.ContentTypeSet("application/json"))
	// r.Use(middleware.RealIP)
	// r.Use(middleware.Recoverer)
	r.Use(packmiddleware.GzipMiddleware)

	r.Head("/", func(rw http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "Content-Type: application/json")
	})
	// GET http://localhost:8080/value/counter/testSetGet163
	// app.HandleRequest()
	r.Get("/", handler.MetricsHandlerFunc)
	r.Get("/ping", handler.PingDB)
	r.Post("/updates", handler.UpdateJSONMetricsHandlerFunc)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handler.UpdateJSONMetricHandlerFunc)
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}", handler.UpdateMetricsHandlerFunc)
			r.Post("/{name}/{value}/", handler.UpdateMetricsHandlerFunc)
			r.Get("/{name}/{value}/", handler.UpdateMetricsHandlerFunc)

		})
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", handler.GetMetricsJSONHandlerFunc)
		r.Get("/", handler.GetMetricHandlerFunc)
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/", handler.GetMetricHandlerFunc)
			r.Get("/{name}", handler.GetMetricHandlerFunc)
			r.Get("/{name}/", handler.GetMetricHandlerFunc)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("./doc.json"), // Ссылка на ваш swagger.json
	))

	// log.Fatal(http.ListenAndServe(flags.Parse(), r))
	return r
}
