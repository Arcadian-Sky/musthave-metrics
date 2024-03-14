package flags

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	serverAddress  = ""
	pollInterval   = time.Second
	reportInterval = time.Second
)

func Parse() {
	end := flag.String("a", "localhost:8080", "endpoint")
	repI := flag.Int("r", 2, "reportInterval")
	polI := flag.Int("p", 10, "pollInterval")
	flag.Parse()

	serverAddress = "http://" + *end
	reportInterval = time.Duration(*repI) * time.Second
	pollInterval = time.Duration(*polI) * time.Second

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverAddress = "http://" + envRunAddr
	}

	if envRepI := os.Getenv("REPORT_INTERVAL"); envRepI != "" {
		durationRepI, err := strconv.Atoi(envRepI)
		if err != nil {
			fmt.Println("Error parsing REPORT_INTERVAL:", err)
			return
		}
		reportInterval = time.Duration(durationRepI) * time.Second
	}

	if envPolI := os.Getenv("POLL_INTERVAL"); envPolI != "" {
		durationPolI, err := strconv.Atoi(envPolI)
		if err != nil {
			fmt.Println("Error parsing POLL_INTERVAL:", err)
			return
		}
		pollInterval = time.Duration(durationPolI) * time.Second
	}
}

func GetServerAddress() string {
	return serverAddress
}

func GetReportInterval() time.Duration {
	return reportInterval
}

func GetPollInterval() time.Duration {
	return pollInterval
}
