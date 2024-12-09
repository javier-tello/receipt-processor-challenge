package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/javier-tello/receipt-processor-challenge/internal/handlers"
	"github.com/javier-tello/receipt-processor-challenge/internal/models"
	"github.com/javier-tello/receipt-processor-challenge/internal/services"
)

func main() {
	// Initialize dependencies
	receiptRepo := models.NewInMemoryReceiptRepo() // Replace with a real repository in production
	receiptService := services.NewUserService(receiptRepo)
	receiptHandler := handlers.NewReceiptHandler(receiptService)

	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", receiptHandler.CreateReceipt)
	r.HandleFunc("/receipts/{id}/points:", receiptHandler.GetReceiptByID)
	r.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	port := ":3000"
	log.Printf("Starting receipt-processor-challenge simple web server on : %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
