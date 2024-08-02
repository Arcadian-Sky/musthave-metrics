package sender

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
)

const UpdatePathOne = "/update"
const UpdatePathPack = "/updates"

type Sender struct {
	getHash       string
	serverAddress string
	cryptoKey     *rsa.PublicKey
}

func NewSender(config *flags.Config) *Sender {
	cKp, ok := config.GetCryptoKeyPath()
	sender := Sender{
		getHash:       config.GetHash(),
		serverAddress: config.GetServerAddress(),
	}
	if ok {
		sender.cryptoKey = cKp
	}
	return &sender
}

// Отправляем запрос на сервер
func (s *Sender) SendMetricJSON(m any, method string) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error marshaling metrics:", err)
		return err
	}
	fmt.Printf("m: %v\n", string(jsonData))
	fmt.Printf("m: %v\n", []byte(jsonData))

	// Формируем адрес запроса
	url := fmt.Sprintf("%s"+method, s.serverAddress)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if s.cryptoKey != nil {
		// Шифруем данные
		encryptedMessage, err := s.encryptMessage([]byte(jsonData), s.cryptoKey)
		if err != nil {
			log.Fatalf("Ошибка при шифровании сообщения: %v", err)
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(encryptedMessage))
		if err != nil {
			return err
		}
	}

	// Создание HTTP-запроса POST
	req.Header.Set("Content-Type", "application/json")

	hashKey := s.getHash
	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(jsonData)
		dst := h.Sum(nil)
		// fmt.Printf("dst: %v\n", hex.EncodeToString(dst))
		req.Header.Set("HashSHA256", hex.EncodeToString(dst))
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//
	// Отправляем запрос на сервер
	// resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	fmt.Printf("Metrics did not sent: \n")
	// 	return err
	// }
	defer resp.Body.Close()

	return nil
}

func (s *Sender) encryptMessage(message []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

func (s *Sender) SendMetricValue(mType string, mName string, mValue interface{}) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Формируем адрес запроса
	url := fmt.Sprintf("%s/update/"+mType+"/%s/%v", s.serverAddress, mName, mValue)

	// Отправляем запрос на сервер
	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		fmt.Printf("Metric did not sent: %s\n", mName)
		return err
	}

	// Печатаем результат отправки (для демонстрации, лучше использовать логгер)
	// fmt.Printf("Metric sent: %s\n", mName)
	defer resp.Body.Close()

	return nil
}
