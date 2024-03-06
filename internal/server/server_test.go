package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arcadian-Sky/musthave-metrics/internal/handler"
	"github.com/stretchr/testify/assert"
)

func Test_methodCheckerMiddleware(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "empty path",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
			request: "/",
		},
		{
			name: "wrong path test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
			request: "/test",
		},
		{
			name: "wrong path counter 1",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
			},
			request: "/update/counter/",
		},
		{
			name: "wrong path counter 2",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
			request: "/update/counter/someName/",
		},
		{
			name: "ok path counter",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
			request: "/update/counter/someName/100/",
		},
		{
			name: "wrong path gauge 1",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
			},
			request: "/update/gauge/",
		},
		{
			name: "wrong path gauge 2",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
			request: "/update/gauge/someName/",
		},
		{
			name: "ok path gauge",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
			request: "/update/gauge/someName/100.001/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Path not allowed", http.StatusBadRequest)
			})
			mux.HandleFunc("/metrics/", handler.MetricsHandler())
			mux.Handle("/update/", сonveyor(
				http.HandlerFunc(handler.UpdateMetricsHandler()),
				methodCheckerMiddleware(),
				сontentTypeCheckerMiddleware("text/plain"),
			))
			request.Body.Close()

			mux.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}

}
