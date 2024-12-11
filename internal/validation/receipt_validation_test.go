package validation

import (
	"log"
	"testing"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
)

func TestReceiptValidator_Validate_EmptyRetailer(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	expectedError := "retailer is required"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("retailer name should not be allowed to be empty")
	}
}

func TestReceiptValidator_Validate_InvalidRetailer(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "T@rget",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	expectedError := "invalid retailer name, must be alphanumeric characters and the following special characters are allow: & and -"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("retailer name should not allow special characters except & and -")
	}
}

func TestReceiptValidator_Validate_EmptyPurchaseDate(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	expectedError := "purchase date is required"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("purchase date is not allowed to be empty")
	}
}

func TestReceiptValidator_Validate_InvalidPurchaseDate(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "12-11-2024",
		PurchaseTime: "13:13",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	expectedError := "invalid purchase date, must be in YYYY-MM-DD format"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Purchase date must be in the format YYYY-MM-DD")
	}
}

func TestReceiptValidator_Validate_EmptyPurchaseTime(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "12-11-2024",
		PurchaseTime: "",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	expectedError := "purchase time is required"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Purchase time cannot be empty")
	}
}

func TestReceiptValidator_Validate_InvalidPurchaseTime(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "5:45pm",
		Total:        "1.25",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	expectedError := "invalid time, must be in HH:MM and in military format"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Purchase time must be in format HH:MM and in military time")
	}
}

func TestReceiptValidator_Validate_EmptyTotal(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "05:45",
		Total:        "",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	expectedError := "total is required"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("total is required")
	}
}

func TestReceiptValidator_Validate_InvalidTotal(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "05:45",
		Total:        "1.111111",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "1.25"}}}

	expectedError := "invalid total, must be in ##.## format"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Purchase total must be in format ##.##")
	}
}

func TestReceiptValidator_Validate_EmptyItems(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "05:45",
		Total:        "1.11",
		Items:        []models.Item{}}

	expectedError := "item(s) are required"
	actualErr := validation.Validate(receipt)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Items cannot be empty")
	}
}

func TestReceiptValidator_Validate_EmptyShortDescription(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "05:45",
		Total:        "1.11",
		Items:        []models.Item{{ShortDescription: "", Price: "1.25"}}}

	expectedError := "short description is required and missing in index 0 of items"
	actualErr := validation.Validate(receipt)
	log.Println(actualErr)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Short description cannot be empty")
	}
}

func TestReceiptValidator_Validate_InvalidShortDescription(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "05:45",
		Total:        "1.11",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz & Coke - 12-oz", Price: "1.25"}}}

	expectedError := "invalid short description in index 0 of items, must be alphanumeric characters with - being the only allowed special character"
	actualErr := validation.Validate(receipt)
	log.Println(actualErr)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Short description can only contain alphanumeric characters and is only allowed a hyphon")
	}
}

func TestReceiptValidator_Validate_EmptyPrice(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "05:45",
		Total:        "1.11",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: ""}}}

	expectedError := "price is required and missing in index 0 of items"
	actualErr := validation.Validate(receipt)
	log.Println(actualErr)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Price cannot be empty")
	}
}

func TestReceiptValidator_Validate_InvalidPrice(t *testing.T) {
	validation := ReceiptValidator{}
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "05:45",
		Total:        "1.11",
		Items:        []models.Item{{ShortDescription: "Pepsi - 12-oz", Price: "2.22222"}}}

	expectedError := "invalid total in index 0 of items, must be in ##.## format"
	actualErr := validation.Validate(receipt)
	log.Println(actualErr)
	if actualErr == nil || actualErr.Error() != expectedError {
		t.Errorf("Price must be in format ##.##")
	}
}
