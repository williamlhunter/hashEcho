package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func computeHash(data []byte, iterations int) string {
	hash := data
	for i := 0; i < iterations; i++ {
		h := sha256.Sum256(hash)
		hash = h[:]
	}
	return hex.EncodeToString(hash)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Default to 1000 iterations if not specified
	iterations := 1000
	if iterStr := r.URL.Query().Get("iterations"); iterStr != "" {
		if iter, err := strconv.Atoi(iterStr); err == nil && iter > 0 {
			iterations = iter
		}
	}

	hash := computeHash(body, iterations)
	response := fmt.Sprintf("%s\n%s\n", string(body), hash)

	log.Printf("%s %s - %d iterations - hash: %s", r.Method, r.URL.Path, iterations, hash[:16])

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "OK\n")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", echoHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("Starting echo server on port %s", port)
	log.Printf("Endpoints:")
	log.Printf("  / - Echo endpoint (echoes all requests)")
	log.Printf("  /health - Health check endpoint")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
