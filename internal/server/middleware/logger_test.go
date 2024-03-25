package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {

}

func Test_loggingResponseWriter_Write(t *testing.T) {
	recorder := httptest.NewRecorder()

	lrw := &loggingResponseWriter{
		ResponseWriter: recorder,
		responseData: &responseData{
			status: 0,
			size:   0,
		},
	}

	data := []byte("test data")
	gotN, err := lrw.Write(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if gotN != len(data) {
		t.Errorf("Unexpected number of bytes written. Expected %d, got %d", len(data), gotN)
	}

	res := recorder.Result()
	defer res.Body.Close()
}

func Test_loggingResponseWriter_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	lrw := &loggingResponseWriter{
		ResponseWriter: recorder,
		responseData: &responseData{
			status: 0,
			size:   0,
		},
	}

	lrw.WriteHeader(http.StatusOK)

	if recorder.Result().StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code. Expected %d, got %d", http.StatusOK, recorder.Result().StatusCode)
	}

	res := recorder.Result()
	defer res.Body.Close()
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	if logger == nil {
		t.Error("Expected non-nil logger, got nil")
	}
}
