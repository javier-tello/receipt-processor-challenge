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

func NewReceiptHandler(receiptService *services.ReceiptService, validator validation.ReceiptValidator) *ReceiptHandler {
	return &ReceiptHandler{
		ReceiptService: receiptService,
		Validator:      validator,
	}
}

func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		log.Printf("Failed to decode receipt JSON: %v", err)
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	if err := h.Validator.ValidateReceipt(receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Processing receipt")
	receiptID := h.ReceiptService.ProcessReceipt(receipt)

	log.Printf("Receipt ID: %s successfully processed", receiptID)
	jsonResponse(w, http.StatusCreated, map[string]string{"id": receiptID})
}

func (h *ReceiptHandler) GetPointsForReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptID := vars["id"]

	if err := h.Validator.ValidateReceiptID(receiptID); err != nil {
		log.Printf("Received invalid receipt id: %s", receiptID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pointsForReceipt, err := h.ReceiptService.CalculateTotalPointsForReceipt(receiptID)

	if err != nil {
		log.Printf("Receipt ID %s not found: %v", receiptID, err)
		http.Error(w, "No receipt found for that ID", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"points": strconv.Itoa(pointsForReceipt)})
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
