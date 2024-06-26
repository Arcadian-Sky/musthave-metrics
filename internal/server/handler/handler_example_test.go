package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

func ExampleHandler() {
	h := &Handler{}

	r := chi.NewRouter()
	r.Get("/", h.MetricsHandlerFunc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
}

func ExampleHandler_GetMetricHandlerFunc() {
	h := &Handler{}

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", h.GetMetricHandlerFunc)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/someName", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
}
