package flags

import (
	"flag"
	"os"
)

func Parse() string {
	end := flag.String("a", ":8080", "endpoint address")
	flag.Parse()
	endpoint := *end

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		endpoint = envRunAddr
	}
	return endpoint
}
