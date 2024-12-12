package validation_test

import (
	"testing"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
	"github.com/javier-tello/receipt-processor-challenge/internal/validation"
)

func TestValidateReceiptID(t *testing.T) {
	validator := &validation.ReceiptValidator{}

	tests := []struct {
		name      string
		receiptID string
		expectErr bool
	}{
		{"Valid Receipt ID", "123456", false},
		{"Empty Receipt ID", "", true},
		{"Whitespace-only Receipt ID", "   ", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateReceiptID(test.receiptID)
			if (err != nil) != test.expectErr {
				t.Errorf("ValidateReceiptID(%q) error = %v, expectErr = %v", test.receiptID, err, test.expectErr)
			}
		})
	}
}

func TestValidateReceipt(t *testing.T) {
	validator := &validation.ReceiptValidator{}

	validItems := []models.Item{
		{ShortDescription: "Item 1", Price: "12.34"},
		{ShortDescription: "Item 2", Price: "45.67"},
	}

	tests := []struct {
		name      string
		receipt   models.Receipt
		expectErr bool
	}{
		{
			name: "Valid Receipt",
			receipt: models.Receipt{
				Retailer:     "Retailer 1",
				PurchaseDate: "2024-12-11",
				PurchaseTime: "14:30",
				Total:        "58.01",
				Items:        validItems,
			},
			expectErr: false,
		},
		{
			name: "Missing Retailer",
			receipt: models.Receipt{
				PurchaseDate: "2024-12-11",
				PurchaseTime: "14:30",
				Total:        "58.01",
				Items:        validItems,
			},
			expectErr: true,
		},
		{
			name: "Invalid Purchase Date",
			receipt: models.Receipt{
				Retailer:     "Retailer 1",
				PurchaseDate: "11-12-2024",
				PurchaseTime: "14:30",
				Total:        "58.01",
				Items:        validItems,
			},
			expectErr: true,
		},
		{
			name: "Invalid Total Format",
			receipt: models.Receipt{
				Retailer:     "Retailer 1",
				PurchaseDate: "2024-12-11",
				PurchaseTime: "14:30",
				Total:        "58.0",
				Items:        validItems,
			},
			expectErr: true,
		},
		{
			name: "Empty Items List",
			receipt: models.Receipt{
				Retailer:     "Retailer 1",
				PurchaseDate: "2024-12-11",
				PurchaseTime: "14:30",
				Total:        "58.01",
				Items:        []models.Item{},
			},
			expectErr: true,
		},
		{
			name: "Invalid Item Short Description",
			receipt: models.Receipt{
				Retailer:     "Retailer 1",
				PurchaseDate: "2024-12-11",
				PurchaseTime: "14:30",
				Total:        "58.01",
				Items: []models.Item{
					{ShortDescription: "", Price: "12.34"},
				},
			},
			expectErr: true,
		},
		{
			name: "Invalid Item Price",
			receipt: models.Receipt{
				Retailer:     "Retailer 1",
				PurchaseDate: "2024-12-11",
				PurchaseTime: "14:30",
				Total:        "58.01",
				Items: []models.Item{
					{ShortDescription: "Item 1", Price: "12.3"},
				},
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validator.ValidateReceipt(test.receipt)
			if (err != nil) != test.expectErr {
				t.Errorf("ValidateReceipt(%+v) error = %v, expectErr = %v", test.receipt, err, test.expectErr)
			}
		})
	}
}

func TestIsValidPurchaseDate(t *testing.T) {
	tests := []struct {
		name         string
		purchaseDate string
		expectErr    bool
	}{
		{"Valid Date", "2024-12-11", true},
		{"Invalid Format", "11-12-2024", false},
		{"Invalid Month", "2024-13-11", false},
		{"Invalid Day", "2024-12-32", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validation.IsValidPurchaseDate(test.purchaseDate)
			if result != test.expectErr {
				t.Errorf("IsValidPurchaseDate(%q) = %v, expected %v", test.purchaseDate, result, test.expectErr)
			}
		})
	}
}
