package sender

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
	pb "github.com/Arcadian-Sky/musthave-metrics/internal/agent/generated/protoagent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const UpdatePathOne = "/update"
const UpdatePathPack = "/updates"

type Sender struct {
	getHash       string
	serverAddress string
	tcpEnabled    bool
	tcpEndpoint   string
	cryptoKey     *rsa.PublicKey
}

func NewSender(config *flags.Config) *Sender {
	cKp, ok := config.GetCryptoKeyPath()
	sender := Sender{
		getHash:       config.GetHash(),
		serverAddress: config.GetServerAddress(),
		tcpEnabled:    config.GetTcpEnable(),
		tcpEndpoint:   config.GetTEndpoint(),
	}
	if ok {
		sender.cryptoKey = cKp
	}
	return &sender
}

// Отправляем запрос на сервер
func (s *Sender) SendMetricJSON(m any, method string) error {
	if s.tcpEnabled {
		return s.SendMetricJSONbyHTTP(m, method)
	} else {
		return s.SendMetricJSONbyGRPC(m, method)
	}
}

func (s *Sender) SendMetricValue(mType string, mName string, mValue interface{}) error {
	if s.tcpEnabled {
		return s.SendValueByGRPC(mType, mName, mValue)
	} else {
		return s.SendValueByHTTP(mType, mName, mValue)
	}
}

func (s *Sender) encryptMessage(message []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

func (s *Sender) getAgentIP() string {
	//  Interfaces returns a list of the system's network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error getting network interfaces: %v\n", err)
		return "unknown"
	}
	// Addrs returns a list of unicast interface addresses for a specific interface
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("Error getting addresses for interface %s: %v\n", iface.Name, err)
			continue
		}
		// Возвращаем первый подходящий IP-адрес
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			fmt.Printf("ipNet.IP: %v\n", ipNet.IP)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return "unknown"
}

func (s *Sender) SendMetricJSONbyHTTP(m any, method string) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error marshaling metrics:", err)
		return err
	}

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
	agentIP := s.getAgentIP() // Получаем IP-адрес агента
	// Создание HTTP-запроса POST
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", agentIP)
	hashKey := s.getHash
	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(jsonData)
		dst := h.Sum(nil)
		req.Header.Set("HashSHA256", hex.EncodeToString(dst))
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *Sender) SendMetricJSONbyGRPC(m any, method string) error {
	// Подключение к gRPC серверу
	conn, err := grpc.NewClient(
		s.tcpEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	// Создание клиента gRPC
	client := pb.NewAgentServiceClient(conn)

	// Сериализация данных
	jsonData, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("error marshaling metrics: %w", err)
	}

	// Создание запроса
	req := pb.MetricJSONRequest{
		JsonString: string(jsonData),
	}

	if s.cryptoKey != nil {
		// Шифруем данные
		encryptedMessage, err := s.encryptMessage(jsonData, s.cryptoKey)
		if err != nil {
			return fmt.Errorf("error encrypting message: %w", err)
		}
		req.JsonString = string(encryptedMessage)
	}

	// Подпись сообщения
	hashKey := s.getHash
	md := metadata.New(map[string]string{})
	if hashKey != "" {
		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(jsonData)
		dst := h.Sum(nil)
		md = metadata.Join(md, metadata.New(map[string]string{
			"HashSHA256": hex.EncodeToString(dst),
		}))
	}

	// Создание контекста с метаданными
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Отправка запроса
	_, err = client.SendMetricJSON(ctx, &req)
	if err != nil {
		return fmt.Errorf("failed to send metric: %w", err)
	}
	return nil
}

func (s *Sender) SendValueByHTTP(mType string, mName string, mValue interface{}) error {
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

	defer resp.Body.Close()

	return nil
}
func (s *Sender) SendValueByGRPC(mType string, mName string, mValue interface{}) error {
	// Создание подключения
	conn, err := grpc.NewClient(
		s.tcpEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	// Создание клиента
	client := pb.NewAgentServiceClient(conn)

	// Преобразование mValue в строку
	valueStr, ok := mValue.(string)
	if !ok {
		return fmt.Errorf("invalid value type: %T", mValue)
	}

	req := &pb.MetricRequest{
		Type:  mType,
		Name:  mName,
		Value: valueStr,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Отправка метрики
	resp, err := client.SendMetric(ctx, req)
	if err != nil {
		return fmt.Errorf("could not send metric: %v", err)
	}

	log.Printf("Metric sent: %s", resp.GetStatus())
	return nil
}
