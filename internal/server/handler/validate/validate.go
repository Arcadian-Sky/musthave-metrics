package validate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

func CheckMetricTypeAndName(mType, mName string) error {
	//Проверяем передачу типа
	if mType == "" {
		return fmt.Errorf("metric type not provided")
		// http.Error(w, "Metric type not provided", http.StatusNotFound)
		// return
	}
	//Проверяем передачу имени
	if mName == "" {
		return fmt.Errorf("metric name not provided")
		// http.Error(w, "Metric name not provided", http.StatusNotFound)
		// return
	}
	return nil
}

func GetHashHead(r *http.Request) string {
	return r.Header.Get("HashSHA256")
}

func CheckHash(sha string, body []byte, key string) error {
	if sha != "" {
		data, err := hex.DecodeString(sha)
		if err != nil {
			return err
		}
		h := hmac.New(sha256.New, []byte(key))
		h.Write(body)
		dst := h.Sum(nil)
		// fmt.Printf("data: %v\n", hmac.Equal(dst, data))
		if !hmac.Equal(dst, data) {
			return fmt.Errorf("hash not valid")
		}
	}
	return nil
}
