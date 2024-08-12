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
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/models"

	// pb "github.com/Arcadian-Sky/musthave-metrics/internal/agent/generated/protoagent"
	pb "github.com/Arcadian-Sky/musthave-metrics/gen/proto/api/metrics/v1"
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
	tcpClient     pb.MetricsServiceClient
}

func NewSender(config *flags.Config) *Sender {
	cKp, ok := config.GetCryptoKeyPath()
	sender := Sender{
		getHash:       config.GetHash(),
		serverAddress: config.GetServerAddress(),
		tcpEnabled:    config.GetTCPEnable(),
		tcpEndpoint:   config.GetTEndpoint(),
	}
	if ok {
		sender.cryptoKey = cKp
	}

	if sender.tcpEnabled {
		// Подключение к gRPC серверу
		conn, err := grpc.NewClient(
			sender.tcpEndpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Printf("failed to connect to server: %s", err)
		}
		defer conn.Close()

		// Создание клиента gRPC
		sender.tcpClient = pb.NewMetricsServiceClient(conn)
	}

	return &sender
}

// Отправляем запрос на сервер
func (s *Sender) SendMetricJSON(m any, method string) error {
	if s.tcpEnabled {
		return s.SendMetricJSONbyGRPC(m, method)
	} else {
		return s.SendMetricJSONbyHTTP(m, method)
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
	// fmt.Printf("m: %v\n", m)
	// fmt.Printf("method: %v\n", method)
	var metricData *pb.Metric
	var metricsData []*pb.Metric
	var req any
	// Создание запроса
	// Определение типа данных и создание соответствующего запроса
	switch v := m.(type) {
	case []models.Metrics:
		// Преобразование среза моделей в срез pb.Metric
		for _, metric := range v {
			metricsData = append(metricsData, convertModelMetricsToRPCMetrics(&metric))
		}
		req = metricsData
		// req = &pb.UpdateJSONMetricsRequest{Metrics: metricsData}
	case models.Metrics:
		// Преобразование одной модели в pb.Metric
		metricData = convertModelMetricsToRPCMetrics(&v)
		req = metricData
		// req = &pb.UpdateJSONMetricRequest{Metric: metricData}
	default:
		return fmt.Errorf("unsupported type %T", v)
	}

	hashKey := s.getHash
	md := metadata.New(map[string]string{})

	if s.cryptoKey != nil && hashKey != "" {
		jsonData, _ := json.Marshal(req)
		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(jsonData)
		dst := h.Sum(nil)
		md = metadata.Join(md, metadata.New(map[string]string{
			"HashSHA256": hex.EncodeToString(dst),
		}))
	}

	// Создание контекста с метаданными
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	switch method {
	case UpdatePathOne:
		_, err := s.tcpClient.UpdateJSONMetric(ctx, &pb.UpdateJSONMetricRequest{Metric: metricData})
		if err != nil {
			return fmt.Errorf("failed to send metric: %w", err)
		}
	case UpdatePathPack:
		_, err := s.tcpClient.UpdateJSONMetrics(ctx, &pb.UpdateJSONMetricsRequest{Metrics: metricsData})
		if err != nil {
			return fmt.Errorf("failed to send metrics: %w", err)
		}
	default:
		return fmt.Errorf("unsupported method %s", method)
	}

	// Отправка запроса
	// _, err = s.tcpClient.UpdateJSONMetric(ctx, req)
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
	// Преобразование mValue в строку
	valueStr, ok := mValue.(string)
	if !ok {
		return fmt.Errorf("invalid value type: %T", mValue)
	}

	req := &pb.UpdateMetricRequest{
		Id:    mName,
		Type:  convertStringToMetricType(mType),
		Value: valueStr,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Отправка метрики
	resp, err := s.tcpClient.UpdateMetric(ctx, req)
	if err != nil {
		return fmt.Errorf("could not send metric: %v", err)
	}
	fmt.Printf("Metric sent: %v\n", resp.GetMetrics())
	return nil
}

func convertStringToMetricType(s string) pb.Type {
	switch s {
	case "gauge":
		return pb.Type_TYPE_GAUGE
	case "counter":
		return pb.Type_TYPE_COUNTER
	default:
		return pb.Type_TYPE_UNSPECIFIED
	}
}

func convertModelMetricsToRPCMetrics(metric *models.Metrics) *pb.Metric {
	var pbMetric pb.Metric

	// Установка идентификатора и типа метрики
	pbMetric.Id = metric.ID
	pbMetric.Type = convertStringToMetricType(metric.MType)

	// Установка значения метрики в зависимости от типа
	switch metric.MType {
	case "counter":
		if metric.Delta != nil {
			pbMetric.Mvalue = &pb.Metric_Delta{
				Delta: *metric.Delta,
			}
		}
	case "gauge":
		if metric.Value != nil {
			pbMetric.Mvalue = &pb.Metric_Value{
				Value: *metric.Value,
			}
		}
	default:
		fmt.Println("Unknown MetricType")
	}

	return &pbMetric
}
