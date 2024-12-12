package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/javier-tello/receipt-processor-challenge/internal/models"
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

func TestReceiptService_PointCalculations(t *testing.T) {
	mockUUIDGenerator := MockUUIDGenerator{}
	repo := NewMockReceiptRepository(mockUUIDGenerator)
	service := NewReceiptService(repo)

	tests := []struct {
		name                string
		receipt             models.Receipt
		expectedTotalPoints int
	}{
		{
			name: "No Points",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:13",
				Total:        "1.01",
				Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}},
			},
			expectedTotalPoints: 0,
		},
		{
			name: "Valid Total Ending in 0",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:13",
				Total:        "1.00",
				Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}},
			},
			expectedTotalPoints: 75,
		},
		{
			name: "Retailer Name",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:13",
				Total:        "1.01",
				Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}},
			},
			expectedTotalPoints: 6,
		},
		{
			name: "Valid Total Ending in a multiple of 25",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:13",
				Total:        "1.75",
				Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}},
			},
			expectedTotalPoints: 25,
		},
		{
			name: "Odd Day Purchase Date",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:13",
				Total:        "1.26",
				Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}},
			},
			expectedTotalPoints: 6,
		},
		{
			name: "Valid Time for Bonus Points",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "14:30",
				Total:        "1.25",
				Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}},
			},
			expectedTotalPoints: 35,
		},
		{
			name: "Valid Short Description Bonus",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:13",
				Total:        "1.01",
				Items:        []models.Item{{ShortDescription: "Emils Cheese Pizza", Price: "12.25"}},
			},
			expectedTotalPoints: 3,
		},
		{
			name: "Multiple Item Pairs",
			receipt: models.Receipt{
				Retailer:     "",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "13:13",
				Total:        "1.01",
				Items: []models.Item{
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
				},
			},
			expectedTotalPoints: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			receiptID, err := service.ProcessReceipt(test.receipt)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			calculatedPoints, err := service.CalculateTotalPointsForReceipt(receiptID)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if test.expectedTotalPoints != calculatedPoints {
				t.Errorf("Test case: %s - Expected points: %d, got: %d", test.name, test.expectedTotalPoints, calculatedPoints)
			}
		})
	}
}

func TestReceiptService_MissingReceiptError(t *testing.T) {
	mockUUIDGenerator := MockUUIDGenerator{}
	repo := NewMockReceiptRepository(mockUUIDGenerator)
	service := NewReceiptService(repo)

	receiptID := "non-existent-id"

	_, err := service.CalculateTotalPointsForReceipt(receiptID)
	if err == nil {
		t.Fatalf("Expected error for missing receipt, got nil")
	}
}
