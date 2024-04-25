package validate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"testing"
)

func TestCheckMetricTypeAndName(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		metricName string
		wantErr    bool
	}{
		{
			name:       "Valid parameters",
			metricType: "cpu",
			metricName: "usage",
			wantErr:    false,
		},
		{
			name:       "Missing metric type",
			metricType: "",
			metricName: "usage",
			wantErr:    true,
		},
		{
			name:       "Missing metric name",
			metricType: "cpu",
			metricName: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckMetricTypeAndName(tt.metricType, tt.metricName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckMetricTypeAndName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetHashHead(t *testing.T) {
	// Создаем заголовок с хэшем
	header := http.Header{}
	header.Set("HashSHA256", "testHash")

	// Создаем запрос с нашим заголовком
	req := &http.Request{Header: header}

	tests := []struct {
		name string
		args *http.Request
		want string
	}{
		{
			name: "HashSHA256 header exists",
			args: req,
			want: "testHash",
		},
		{
			name: "HashSHA256 header does not exist",
			args: &http.Request{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHashHead(tt.args); got != tt.want {
				t.Errorf("GetHashHead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckHash(t *testing.T) {
	// Устанавливаем тестовые данные
	body := []byte("test body")
	key := "test key"
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write(body)
	expectedHash := hex.EncodeToString(hash.Sum(nil))

	tests := []struct {
		name    string
		sha     string
		body    []byte
		key     string
		wantErr bool
	}{
		{
			name:    "Valid hash",
			sha:     expectedHash,
			body:    body,
			key:     key,
			wantErr: false,
		},
		{
			name:    "Invalid hash",
			sha:     "invalidHash",
			body:    body,
			key:     key,
			wantErr: true,
		},
		{
			name:    "Empty hash",
			sha:     "",
			body:    body,
			key:     key,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckHash(tt.sha, tt.body, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
