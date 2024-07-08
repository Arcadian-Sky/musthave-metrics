package main

import (
	"fmt"

	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/controller"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	config, err := flags.Parse()
	if err != nil {
		panic("Panic error in flags parsing")
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	controller.NewCollectAndSendMetricsService(&config).Run()
}
