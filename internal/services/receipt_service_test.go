package services

import (
	"testing"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
)

type MockReceiptRepository struct {
	receipts  map[int]models.Receipt
	idCounter int
}

func NewMockReceiptRepository() *MockReceiptRepository {
	return &MockReceiptRepository{receipts: make(map[int]models.Receipt), idCounter: 0}
}

func (m *MockReceiptRepository) ProcessReceipt(receipt models.Receipt) int {
	receipt.ID = m.idCounter
	m.idCounter++
	m.receipts[receipt.ID] = receipt

	return receipt.ID
}

func (m *MockReceiptRepository) FindByID(receiptID int) (models.Receipt, bool) {
	receipt, exists := m.receipts[receiptID]

	return receipt, exists
}

func TestReceiptService_CalculateTotalPointsForReceipt(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 31
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}
