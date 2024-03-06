package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/handler"
)

func Test_methodCheckerMiddleware(t *testing.T) {
	tests := []struct {
		name string
		want Middleware
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := methodCheckerMiddleware(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("methodCheckerMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.MetricsHandler())

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status":"metrics"}`

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestContentTypeCheckerMiddleware(t *testing.T) {
	handlerToWrap := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := server.сontentTypeCheckerMiddleware("application/json")(handlerToWrap)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("handler returned unexpected content type: got %v want %v",
			contentType, expected)
	}
}
