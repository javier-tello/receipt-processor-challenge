package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/javier-tello/receipt-processor-challenge/internal/handlers"
	"github.com/javier-tello/receipt-processor-challenge/internal/repositories"
	"github.com/javier-tello/receipt-processor-challenge/internal/services"
)

func main() {
	receiptRepo := repositories.NewInMemoryReceiptRepo()
	receiptService := services.NewReceiptService(receiptRepo)
	receiptHandler := handlers.NewReceiptHandler(receiptService)

	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", receiptHandler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", receiptHandler.GetReceiptByID).Methods("GET")

	port := ":3000"
	log.Printf("Starting receipt-processor-challenge simple web server on : %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
