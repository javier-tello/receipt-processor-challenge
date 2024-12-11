package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
	"github.com/javier-tello/receipt-processor-challenge/internal/repositories"
	"github.com/javier-tello/receipt-processor-challenge/internal/services"
)

func TestHandler_ProcessReceipt(t *testing.T) {
	repo := repositories.NewInMemoryReceiptRepo()
	service := services.NewReceiptService(repo)
	handler := NewReceiptHandler(service)

	payload := `{"retailer": "Target", "purchaseDate": "2022-01-02","purchaseTime": "13:13", "total": "1.25", "items": [{"shortDescription": "Pepsi - 12-oz", "price": "1.25"}]}`
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ProcessReceipt(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	expectedResponse := `{"id":"0"}`
	actualResponse := rec.Body.String()

	var expected, actual map[string]interface{}
	if err := json.Unmarshal([]byte(expectedResponse), &expected); err != nil {
		t.Fatalf("Failed to parse expected JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(actualResponse), &actual); err != nil {
		t.Fatalf("Failed to parse actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected response: %v, got: %v", expected, actual)
	}
}

func TestHandler_GetPointsForReceipt(t *testing.T) {
	repo := repositories.NewInMemoryReceiptRepo()
	service := services.NewReceiptService(repo)
	handler := NewReceiptHandler(service)

	receipt := repo.ProcessReceipt(models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})

	receiptID := strconv.Itoa(receipt)

	req := httptest.NewRequest(http.MethodGet, "/receipts/"+receiptID+"/points", nil)
	rec := httptest.NewRecorder()

	log.Println(req.RequestURI)

	router := mux.NewRouter()
	router.HandleFunc("/receipts/{id}/points", handler.GetPointsForReceipt).Methods("GET")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedResponse := `{"points":"31"}`
	actualResponse := rec.Body.String()

	var expected, actual map[string]interface{}
	if err := json.Unmarshal([]byte(expectedResponse), &expected); err != nil {
		t.Fatalf("Failed to parse expected JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(actualResponse), &actual); err != nil {
		t.Fatalf("Failed to parse actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected response: %v, got: %v", expected, actual)
	}

}
