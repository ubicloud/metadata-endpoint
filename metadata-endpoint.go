package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/time/rate"
)

// Rate limiter allowing 10 requests per second
var limiter = rate.NewLimiter(10, 4)

func fileHandler(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit the request size
		r.Body = http.MaxBytesReader(w, r.Body, 1024)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Read the specified file
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Write the file content to the response
		w.Write(data)
	}
}

func mainInternal(listenAndServe func(addr string, handler http.Handler) error) (int, error) {
	vmInhostName := os.Getenv("VM_INHOST_NAME")
	if vmInhostName == "" {
		return 1, fmt.Errorf("environment variable VM_INHOST_NAME is not set")
	}

	ipv6Address := os.Getenv("IPV6_ADDRESS")
	if ipv6Address == "" {
		return 1, fmt.Errorf("environment variable IPV6_ADDRESS is not set")
	}

	certFilePath := fmt.Sprintf("/vm/%s/cert/cert.pem", vmInhostName)
	keyFilePath := fmt.Sprintf("/vm/%s/cert/key.pem", vmInhostName)

	http.HandleFunc("/load-balancer/cert.pem", fileHandler(certFilePath))
	http.HandleFunc("/load-balancer/key.pem", fileHandler(keyFilePath))

	address := fmt.Sprintf("[%s]:8080", ipv6Address) // Format for IPv6

	log.Printf("Server starting on %s...\n", address)
	return 0, listenAndServe(address, nil)
}
