package middleware

import (
	"net"
	"net/http"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
)

func SubnetMiddleware(c flags.InitedFlags) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.TrustedSubnetS != "" {
				realIP := r.Header.Get("X-Real-IP")
				if realIP == "" {
					http.Error(w, "X-Real-IP header is missing", http.StatusBadRequest)
					return
				}

				ip := net.ParseIP(realIP)
				if ip == nil {
					http.Error(w, "Invalid X-Real-IP address", http.StatusBadRequest)
					return
				}

				if !c.TrustedSubnet.Contains(ip) {
					http.Error(w, "Forbidden: IP address not in trusted subnet", http.StatusForbidden)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
