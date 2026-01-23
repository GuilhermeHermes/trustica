// Package main provides a simple HTTPS server for testing certificate trust.
// This server uses a certificate signed by our custom CA.
// Tools that don't trust our CA will fail to connect.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	certFile := getEnv("TLS_CERT", "/certs/server.pem")
	keyFile := getEnv("TLS_KEY", "/certs/server-key.pem")
	addr := getEnv("ADDR", ":8443")

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok", "message": "Trustica test server"}`)
	})

	// Endpoint for testing
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"healthy": true}`)
	})

	// Simulated git repository info (for git ls-remote test)
	mux.HandleFunc("/info/refs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
		w.WriteHeader(http.StatusOK)
		// Minimal git protocol response
		fmt.Fprintf(w, "001e# service=git-upload-pack\n0000")
	})

	// Simulated PyPI simple index (for pip test)
	mux.HandleFunc("/simple/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<!DOCTYPE html><html><body><h1>Simple Index</h1></body></html>`)
	})

	log.Printf("Starting HTTPS server on %s", addr)
	log.Printf("Using cert: %s", certFile)
	log.Printf("Using key:  %s", keyFile)

	err := http.ListenAndServeTLS(addr, certFile, keyFile, mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
