package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
	"github.com/javier-tello/receipt-processor-challenge/internal/services"
	"github.com/javier-tello/receipt-processor-challenge/internal/validation"
)

type ReceiptHandler struct {
	ReceiptService *services.ReceiptService
	Validator      *validation.ReceiptValidator
}

func NewReceiptHandler(receiptService *services.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{ReceiptService: receiptService}
}

func (h *ReceiptHandler) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.Validator.Validate(receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.ReceiptService.CreateReceipt(receipt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetUserByID handles GET /users/{id}.
func (h *ReceiptHandler) GetReceiptByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptID := vars["id"]
	convertedReceiptId, err := strconv.Atoi(receiptID)
	if err != nil || convertedReceiptId < 0 {
		http.Error(w, "Invalid receipt ID", http.StatusBadRequest)
		return
	}

	// Use the service to fetch the user
	response, err := h.ReceiptService.GetReceiptByID(convertedReceiptId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the user data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponseData)
}
