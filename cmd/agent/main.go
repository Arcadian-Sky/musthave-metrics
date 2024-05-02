package main

import (
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/controller"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
)

func main() {
	config, err := flags.Parse()
	if err != nil {
		panic("Panic error in flags parsing")
	}
	controller.NewCollectAndSendMetricsService(&config).Run()
}
