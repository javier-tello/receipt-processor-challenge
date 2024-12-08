package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	port := "3000"
	log.Printf("Starting receipt-processor-challenge simple web server on : %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
