package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
	"github.com/javier-tello/receipt-processor-challenge/internal/services"
	"github.com/javier-tello/receipt-processor-challenge/internal/validation"
)

type ReceiptHandler struct {
	ReceiptService *services.ReceiptService
	Validator      validation.ReceiptValidator
}

func NewReceiptHandler(receiptService *services.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{ReceiptService: receiptService}
}

func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.Validator.Validate(receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Processing reciept")
	receiptID, err := h.ReceiptService.ProcessReceipt(receipt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Println("Something went wrong with processing receipt")
		return
	}
	log.Println("Receipt succesfully processed")

	response := map[string]string{
		"id": strconv.Itoa(receiptID),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ReceiptHandler) GetPointsForReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptID := vars["id"]
	convertedReceiptId, err := strconv.Atoi(receiptID)

	if err != nil || convertedReceiptId < 0 {
		http.Error(w, "Invalid receipt ID", http.StatusBadRequest)
		return
	}

	pointsForReceipt, err := h.ReceiptService.CalculateTotalPointsForReceipt(convertedReceiptId)

	if err != nil {
		log.Println("Receipt not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]string{
		"points": strconv.Itoa(pointsForReceipt),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
