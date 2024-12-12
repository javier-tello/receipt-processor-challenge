package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
	"github.com/javier-tello/receipt-processor-challenge/internal/services"
	"github.com/javier-tello/receipt-processor-challenge/internal/validation"
)

type MockUUIDGenerator struct{}

func (m MockUUIDGenerator) New() uuid.UUID {
	return uuid.MustParse("123e4567-e89b-12d3-a456-426614174000") // Fixed UUID for deterministic tests
}

type MockReceiptRepository struct {
	receipts    map[string]models.Receipt
	idGenerator MockUUIDGenerator
}

func NewMockReceiptRepository(generator MockUUIDGenerator) *MockReceiptRepository {
	return &MockReceiptRepository{receipts: make(map[string]models.Receipt), idGenerator: generator}
}

func (m *MockReceiptRepository) ProcessReceipt(receipt models.Receipt) string {
	receiptID := m.idGenerator.New().String()
	m.receipts[receiptID] = receipt
	return receiptID
}

func (m *MockReceiptRepository) FindByID(receiptID string) (models.Receipt, bool) {
	receipt, exists := m.receipts[receiptID]
	return receipt, exists
}

// Helper to set up handler and dependencies
func setupHandler() *ReceiptHandler {
	receiptValidator := validation.ReceiptValidator{}
	mockUUIDGenerator := MockUUIDGenerator{}
	repo := NewMockReceiptRepository(mockUUIDGenerator)
	service := services.NewReceiptService(repo)

	return NewReceiptHandler(service, receiptValidator)
}

// Helper to configure router
func setupRouter(handler *ReceiptHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", handler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", handler.GetPointsForReceipt).Methods("GET")

	return router
}

func TestHandler_ProcessReceipt_ValidPayload(t *testing.T) {
	handler := setupHandler()

	payload := `{"retailer": "Target", "purchaseDate": "2022-01-02", "purchaseTime": "13:13", "total": "1.25", "items": [{"shortDescription": "Pepsi - 12-oz", "price": "1.25"}]}`
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router := setupRouter(handler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	expectedResponse := map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"}
	var actualResponse map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("Failed to parse actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedResponse, actualResponse) {
		t.Errorf("Expected response: %v, got: %v", expectedResponse, actualResponse)
	}
}

func TestHandler_GetPointsForReceipt_ValidID(t *testing.T) {
	handler := setupHandler()

	// Preload a receipt in the repository
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}},
	}
	receiptID, err := handler.ReceiptService.ProcessReceipt(receipt)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/receipts/"+receiptID+"/points", nil)
	rec := httptest.NewRecorder()

	router := setupRouter(handler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedResponse := map[string]interface{}{"points": "31"}
	var actualResponse map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("Failed to parse actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedResponse, actualResponse) {
		t.Errorf("Expected response: %v, got: %v", expectedResponse, actualResponse)
	}
}

func TestHandler_GetPointsForReceipt_InvalidID(t *testing.T) {
	handler := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/receipts/invalid-id/points", nil)
	rec := httptest.NewRecorder()

	router := setupRouter(handler)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
