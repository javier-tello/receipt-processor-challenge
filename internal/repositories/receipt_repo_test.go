package repositories

import (
	"reflect"
	"testing"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
)

func TestInMemoryReceiptRepo_ProcessReceipt(t *testing.T) {
	repo := NewInMemoryReceiptRepo()

	receipt := repo.ProcessReceipt(models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}})

	if receipt != 0 {
		t.Errorf("Expected ID to be '0', got '%d'", receipt)
	}
}

func TestInMemoryReceiptRepo_FindByID(t *testing.T) {
	repo := NewInMemoryReceiptRepo()

	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	receiptID := repo.ProcessReceipt(receipt)
	foundReceipt, exists := repo.FindByID(receiptID)
	if !exists {
		t.Fatalf("User with ID '%d' not found", receiptID)
	}

	if !reflect.DeepEqual(foundReceipt, receipt) {
		t.Errorf("Expected user: %+v, got: %+v", receipt, foundReceipt)
	}
}
