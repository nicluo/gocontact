package gocontact

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func RequestIP(req *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	// For reverse proxy
	forward := req.Header.Get("X-Forwarded-For")

	log.Printf("ip: %s\n", ip)
	log.Printf("forwarded for: %s\n", forward)

	if forward != "" {
		return forward, nil
	}

	return ip, nil
}
