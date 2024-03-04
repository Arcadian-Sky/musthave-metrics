package server

import (
	"fmt"
	"net/http"

	"github.com/Arcadian-Sky/musthave-metrics/internal/handler"
)

type Middleware func(http.Handler) http.Handler

// сontentTypeCheckerMiddleware возвращает middleware, которое проверяет тип данных и устанавливает тип для ответа
func сontentTypeCheckerMiddleware(expectedContentType string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем значение Content-Type из заголовка
			contentType := r.Header.Get("Content-Type")

			// Проверяем, соответствует ли Content-Type ожидаемому значению
			if contentType != expectedContentType {
				http.Error(w, "Error in Content-Type", http.StatusBadRequest)
				return
			}

			// Устанавливаем Content-Type для ответа
			w.Header().Set("Content-Type", expectedContentType)

			// Вызываем следующий обработчик в цепочке
			next.ServeHTTP(w, r)
		})
	}
}

// methodCheckerMiddleware возвращает middleware, которое проверяет метод
func methodCheckerMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Вызываем следующий обработчик в цепочке
			next.ServeHTTP(w, r)
		})
	}
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Привет"))
}

func сonveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func InitServer() {
	// Регистрируем обработчики
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Path not allowed", http.StatusBadRequest)
	})
	http.HandleFunc("/metrics/", handler.MetricsHandler())
	http.Handle("/update/", сonveyor(
		http.HandlerFunc(handler.UpdateMetricsHandler()),
		methodCheckerMiddleware(),
		сontentTypeCheckerMiddleware("text/plain"),
	))

	// Запускаем сервер на порту 8080
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
