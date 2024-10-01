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

func main() {
	vmInhostName := os.Getenv("VM_INHOST_NAME")
	if vmInhostName == "" {
		log.Fatal("Environment variable vm_path is not set")
	}

	ipv6Address := os.Getenv("IPV6_ADDRESS")
	if ipv6Address == "" {
		log.Fatal("Environment variable ipv6_address is not set")
	}

	certFilePath := fmt.Sprintf("/vm/%s/cert/cert.pem", vmInhostName)
	keyFilePath := fmt.Sprintf("/vm/%s/cert/key.pem", vmInhostName)

	http.HandleFunc("/load-balancer/cert.pem", fileHandler(certFilePath))
	http.HandleFunc("/load-balancer/key.pem", fileHandler(keyFilePath))

	address := fmt.Sprintf("[%s]:8080", ipv6Address) // Format for IPv6

	log.Printf("Server starting on %s...\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
