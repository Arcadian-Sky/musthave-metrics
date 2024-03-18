package validate

import (
	"fmt"
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
