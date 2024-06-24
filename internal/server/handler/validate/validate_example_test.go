package validate

import "fmt"

func ExampleCheckMetricTypeAndName() {
	error := CheckMetricTypeAndName("", "")
	fmt.Println(error.Error())

	error = CheckMetricTypeAndName("111", "222")
	fmt.Println(error.Error())

	error = CheckMetricTypeAndName("gauge", "222")
	fmt.Println(error.Error())
}
