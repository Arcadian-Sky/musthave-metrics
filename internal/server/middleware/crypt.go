package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
)

func DecryptMiddleware(c flags.InitedFlags) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			privateKey, _ := c.GetCryptoKey()
			if privateKey != nil {
				// Читаем зашифрованные данные из тела запроса
				encryptedData, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Ошибка при чтении данных", http.StatusBadRequest)
					return
				}
				defer r.Body.Close()

				// Расшифровываем данные
				decryptedData, err := decryptMessage(encryptedData, privateKey)
				if err != nil {
					http.Error(w, "Ошибка при расшифровке данных", http.StatusInternalServerError)
					return
				}

				// Подменяем тело запроса на расшифрованные данные
				r.Body = io.NopCloser(bytes.NewReader(decryptedData))
			}

			// Передаем управление следующему обработчику
			h.ServeHTTP(w, r)
		})
	}
}

func decryptMessage(encryptedMessage []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedMessage)
}
