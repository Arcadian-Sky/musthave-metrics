package noosexit

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMain(t *testing.T) {
	var testdata string
	var err error
	testdata, err = filepath.Abs("testdata/noosexit")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(testdata)
	analysistest.Run(t, testdata, NoOsExitAnalyzer, "./...")

	testdata, err = filepath.Abs("testdata/func")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(testdata)
	analysistest.Run(t, testdata, NoOsExitAnalyzer, "./...")

	testdata, err = filepath.Abs("testdata/main")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(testdata)

	res := analysistest.Run(t, testdata, NoOsExitAnalyzer, "./...")
	fmt.Println(res)

}
