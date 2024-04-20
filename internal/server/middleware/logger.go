package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
		data   string
	}

	// responseWriter - обертка вокруг http.ResponseWriter, чтобы отслеживать код статуса и размер ответа.
	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

// записываем ответ, используя оригинальный http.ResponseWriter
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	if r.responseData.status == 0 {
		r.responseData.status = http.StatusOK
	}
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	// r.responseData.data = string(b)
	return size, err
}

// записываем код статуса, используя оригинальный http.ResponseWriter
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	if r.responseData.status == 0 {
		r.ResponseWriter.WriteHeader(statusCode)
		r.responseData.status = statusCode // захватываем код статуса
	}
}

// getLogger возвращает экземпляр логгера.
func GetLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

// LoggingMiddleware добавляет логирование для каждого запроса и ответа.
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Захватываем все действия, которые выполняются после обработки запроса
		// и передачи управления следующему обработчику.
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData: &responseData{
				status: 0,
				size:   0,
				data:   "",
			},
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger := GetLogger()

		logger.Info("Request processed",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", lw.responseData.status),
			zap.Duration("duration", duration),
			zap.Int("response_size", lw.responseData.size),
			zap.String("response", lw.responseData.data),
		)

		// logger.Info(
		// 	"uri", r.RequestURI,
		// 	"method", r.Method,
		// 	"status", responseData.status,
		// 	"duration", duration,
		// 	"size", responseData.size,
		// )
	})
}
