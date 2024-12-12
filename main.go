package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/javier-tello/receipt-processor-challenge/internal/handlers"
	"github.com/javier-tello/receipt-processor-challenge/internal/repositories"
	"github.com/javier-tello/receipt-processor-challenge/internal/services"
	"github.com/javier-tello/receipt-processor-challenge/internal/validation"
)

func main() {
	receiptValidator := validation.ReceiptValidator{}
	receiptRepo := repositories.NewInMemoryReceiptRepo(nil)
	receiptService := services.NewReceiptService(receiptRepo)
	receiptHandler := handlers.NewReceiptHandler(receiptService, receiptValidator)

	router := setupRouter(receiptHandler)

	port := ":3000"
	log.Printf("Starting receipt-processor-challenge simple web server on : %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}

func setupRouter(handler *handlers.ReceiptHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", handler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", handler.GetPointsForReceipt).Methods("GET")

	return router
}
