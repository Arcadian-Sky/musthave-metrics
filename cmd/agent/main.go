package main

import (
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/controller"
	"github.com/Arcadian-Sky/musthave-metrics/internal/agent/flags"
)

func main() {
	flags.Parse()
	controller.CollectAndSendMetrics()
}
