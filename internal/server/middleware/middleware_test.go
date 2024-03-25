package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentTypeSet(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "Test with text/plain content type",
			args: "text/plain",
			want: "text/plain",
		},
		{
			name: "Test with application/json content type",
			args: "application/json",
			want: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := ContentTypeSet(tt.args)
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				contentType := w.Header().Get("Content-Type")
				if contentType != tt.want {
					t.Errorf("ContentTypeChecker() Content-Type = %v, want %v", contentType, tt.want)
				}
			})
			handler := middleware(testHandler)
			req := httptest.NewRequest("GET", "/", nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)
		})
	}
}
