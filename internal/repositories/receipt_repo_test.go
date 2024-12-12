package repositories

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/javier-tello/receipt-processor-challenge/internal/models"
)

type MockUUIDGenerator struct{}

func (m MockUUIDGenerator) New() uuid.UUID {
	return uuid.MustParse("123e4567-e89b-12d3-a456-426614174000") // Fixed UUID for deterministic tests
}

func TestInMemoryReceiptRepo_ProcessReceipt(t *testing.T) {
	// Arrange: Use the mock UUID generator
	mockGenerator := MockUUIDGenerator{}
	repo := NewInMemoryReceiptRepo(mockGenerator)

	receiptID := repo.ProcessReceipt(models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})

	expectedUUID := "123e4567-e89b-12d3-a456-426614174000"
	if receiptID != expectedUUID {
		t.Errorf("Expected receipt ID to be '%s', got '%s'", expectedUUID, receiptID)
	}

	if _, exists := repo.receipts[expectedUUID]; !exists {
		t.Errorf("Expected receipt to be stored with UUID '%s'", expectedUUID)
	}
}

func TestInMemoryReceiptRepo_FindByID(t *testing.T) {
	repo := NewInMemoryReceiptRepo(nil)

	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	receiptID := repo.ProcessReceipt(receipt)
	foundReceipt, exists := repo.FindByID(receiptID)
	if !exists {
		t.Fatalf("User with ID '%s' not found", receiptID)
	}

	if !reflect.DeepEqual(foundReceipt, receipt) {
		t.Errorf("Expected user: %+v, got: %+v", receipt, foundReceipt)
	}
}
