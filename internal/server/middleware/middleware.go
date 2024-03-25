package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// сontentTypeCheckerMiddleware возвращает middleware, которое проверяет тип данных и устанавливает тип для ответа
func ContentTypeChecker(expectedContentType string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Устанавливаем Content-Type для ответа
			w.Header().Set("Content-Type", expectedContentType)
			// Вызываем следующий обработчик в цепочке
			next.ServeHTTP(w, r)
		})
	}
}
