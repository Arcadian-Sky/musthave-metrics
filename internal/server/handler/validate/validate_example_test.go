package validate

import "fmt"

// ExampleCheckMetricTypeAndName - example test
func ExampleCheckMetricTypeAndName() {
	err := CheckMetricTypeAndName("", "metricName")
	if err != nil {
		fmt.Println(err)
	}

	err = CheckMetricTypeAndName("metricType", "")
	if err != nil {
		fmt.Println(err)
	}

	err = CheckMetricTypeAndName("metricType", "metricName")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// metric type not provided
	// metric name not provided
}
