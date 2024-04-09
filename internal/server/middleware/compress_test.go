package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipMiddleware(t *testing.T) {
	// Создаем запрос с заголовком Content-Encoding: gzip
	reqBody := strings.NewReader("some compressed data")
	req, err := http.NewRequest("POST", "/", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем HTTP запись для записи ответа
	rr := httptest.NewRecorder()

	// Запускаем Middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Set("Content-Encoding", "gzip")
		// Проверяем, что тело запроса было декомпрессировано
		if r.Body == nil {
			t.Error("Request body is nil after decompression")
		}
	})

	// Применяем Middleware к обработчику и выполняем запрос
	GzipMiddleware(handler).ServeHTTP(rr, req)

	// Проверяем статус код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGzipMiddlewareWithGzipRequestBody(t *testing.T) {
	// Создаем новый хендлер, который будет обернут Middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что Middleware правильно принимает запросы с сжатым телом
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}
		defer r.Body.Close()

		// Проверяем, что тело запроса содержит ожидаемые данные
		expected := "Hello, world!"
		if string(body) != expected {
			t.Errorf("Request body = %v, want %v", string(body), expected)
		}
	})

	// Создаем новый запрос, используя созданный хендлер и Middleware
	body := compressString("Hello, world!")
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Encoding", "gzip")

	// Создаем записывающий ResponseRecorder (реализация http.ResponseWriter), чтобы записать ответ
	rr := httptest.NewRecorder()

	// Запускаем Middleware
	GzipMiddleware(handler).ServeHTTP(rr, req)
}

func TestWriteHeader(t *testing.T) {
	// Создаем новый compressWriter
	cw := &compressWriter{
		w:  httptest.NewRecorder(),
		zw: gzip.NewWriter(httptest.NewRecorder()),
	}

	// Вызываем WriteHeader для compressWriter
	cw.WriteHeader(http.StatusOK)

	// Проверяем, что заголовки были установлены корректно
	if got := cw.w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Errorf("Content-Encoding = %v, want %v", got, "gzip")
	}
}

func compressString(s string) *bytes.Buffer {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(s)); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return &b
}
