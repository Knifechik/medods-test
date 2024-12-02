package http

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// getIP gets user ip by request
func getIP(r *http.Request) (string, error) {

	ipAddress := r.RemoteAddr
	fwdAddress := r.Header.Get("X-Forwarded-For")

	if fwdAddress != "" {
		ipAddress = fwdAddress

		ips := strings.Split(fwdAddress, ",")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	}

	if strings.Contains(ipAddress, ":") {
		host, _, err := net.SplitHostPort(ipAddress)
		if err != nil {
			return "", fmt.Errorf("SplitHostPort: %w", err)
		}

		ipAddress = host
	}

	IP := net.ParseIP(ipAddress)

	return IP.String(), nil
}
