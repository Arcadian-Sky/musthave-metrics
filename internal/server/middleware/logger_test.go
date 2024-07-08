package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TesLoggingResponseWriter_GetLogger(t *testing.T) {
	logger := GetLogger()
	if logger == nil {
		t.Error("Expected non-nil logger, got nil")
	}
}

// Тест для метода Write
func TesLoggingResponseWriter_tWrite(t *testing.T) {
	// Создаем фальшивый ResponseWriter
	rec := httptest.NewRecorder()

	// Создаем loggingResponseWriter, оборачивающий фальшивый ResponseWriter
	lrw := &loggingResponseWriter{
		ResponseWriter: rec,
		responseData:   &responseData{},
	}

	// Данные, которые будем писать
	data := []byte("Hello, World!")

	// Вызываем метод Write
	size, err := lrw.Write(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Проверяем размер записанных данных
	if size != len(data) {
		t.Errorf("expected size %d, got %d", len(data), size)
	}

	// Проверяем захваченный размер
	if lrw.responseData.size != len(data) {
		t.Errorf("expected captured size %d, got %d", len(data), lrw.responseData.size)
	}

	// Проверяем статус ответа
	if lrw.responseData.status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, lrw.responseData.status)
	}

	// Проверяем записанные данные
	if rec.Body.String() != string(data) {
		t.Errorf("expected body %s, got %s", string(data), rec.Body.String())
	}
}

// MockLogger используется для проверки того, что нужные данные были записаны в лог
// MockLogger is a mock implementation of zapcore.Core and zapcore.WriteSyncer
type MockLogger struct {
	zapcore.Core
	zapcore.WriteSyncer
}

func (m *MockLogger) Enabled(zapcore.Level) bool {
	return true
}

func (m *MockLogger) With([]zapcore.Field) zapcore.Core {
	return m
}

func (m *MockLogger) Check(zapcore.Entry, *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return nil
}

func (m *MockLogger) Write(e zapcore.Entry, fields []zapcore.Field) error {
	return nil
}

func (m *MockLogger) Sync() error {
	return nil
}

var logger *zap.Logger

// SetLogger устанавливает глобальный логгер
func SetHandlerLogger(l *zap.Logger) {
	logger = l
}

// GetLogger возвращает текущий глобальный логгер
func GetHandlerLogger() *zap.Logger {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return logger
}

func TestHandler_Logger(t *testing.T) {
	mockCore := &MockLogger{}
	mockLogger := zap.New(mockCore)

	SetHandlerLogger(mockLogger)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	loggerHandler := Logger(handler)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	loggerHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем, что ответное тело корректное
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}
