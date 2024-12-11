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

func TestReceiptService_CalculatePoitntsForRetailer(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.26",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 6
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForPurchaseDateOddDay(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:13",
		Total:        "1.26",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 6
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForPurchaseDateEvenDay(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.26",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 0
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForPurchaseTimeNoPoints(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.26",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 0
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForPurchaseTime(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "14:13",
		Total:        "1.26",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 10
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForTotalEndingIn25Multiple(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 25
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForTotalEndingIn0(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.00",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 75
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForItemPairsNoPoints(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.01",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 0
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForItemPairs(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.01",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}, {ShortDescription: "Pepsi - 12-oz", Price: "1.25"}, {ShortDescription: "Pepsi - 12-oz", Price: "1.25"}, {ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 10
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForShortDescriptionNoPoints(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.01",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 0
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}

func TestReceiptService_CalculatePointsForShortDescription(t *testing.T) {
	repo := NewMockReceiptRepository()
	service := NewReceiptService(repo)

	receipt, err := service.ProcessReceipt(models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.01",
		Items:        []models.Item{{ShortDescription: "Emils Cheese Pizza", Price: "12.25"}}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotalPoints := 3
	calculatedTotalPoints, err := service.CalculateTotalPointsForReceipt(receipt)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedTotalPoints != calculatedTotalPoints {
		t.Errorf("Expected points to total: %+d, got: %+d", expectedTotalPoints, calculatedTotalPoints)
	}
}
