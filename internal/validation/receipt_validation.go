package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
)

var (
	receiptIDRegex            = regexp.MustCompile(`^\S+$`)
	retailerRegex             = regexp.MustCompile(`^[\w\s\-&]+$`)
	purchaseTimeRegex         = regexp.MustCompile(`^(2[0-3]|[01][0-9]):[0-5][0-9]$`)
	amountRegex               = regexp.MustCompile(`^\d+\.\d{2}$`)
	itemShortDescriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
)

type ReceiptValidator struct{}

func (uv *ReceiptValidator) ValidateReceiptID(receiptID string) error {
	if strings.TrimSpace(receiptID) == "" {
		return errors.New("please pass in a non-empty id")
	}
	if !receiptIDRegex.MatchString(receiptID) {
		return errors.New("invalid id passed in")
	}
	return nil
}

func (uv *ReceiptValidator) ValidateReceipt(receipt models.Receipt) error {
	var validationErrors []string

	if receipt.Retailer == "" {
		validationErrors = append(validationErrors, "retailer is required")
	}
	if receipt.PurchaseDate == "" {
		validationErrors = append(validationErrors, "purchase date is required")
	}
	if receipt.PurchaseTime == "" {
		validationErrors = append(validationErrors, "purchase time is required")
	}
	if receipt.Total == "" {
		validationErrors = append(validationErrors, "total is required")
	}
	if len(receipt.Items) == 0 {
		validationErrors = append(validationErrors, "item(s) are required")
	}

	if !retailerRegex.MatchString(receipt.Retailer) {
		validationErrors = append(validationErrors, "invalid retailer name")
	}
	if !IsValidPurchaseDate(receipt.PurchaseDate) {
		validationErrors = append(validationErrors, "invalid purchase date, must be in YYYY-MM-DD format")
	}
	if !purchaseTimeRegex.MatchString(receipt.PurchaseTime) {
		validationErrors = append(validationErrors, "invalid purchase time, must be in HH:MM format")
	}
	if !amountRegex.MatchString(receipt.Total) {
		validationErrors = append(validationErrors, "invalid total, must be in ##.## format")
	}

	for i, item := range receipt.Items {
		if item.Price == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("missing price in item at index %d", i))
		}
		if item.ShortDescription == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("missing short description in item at index %d", i))
		}
		if !amountRegex.MatchString(item.Price) {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid price in item at index %d", i))
		}
		if !itemShortDescriptionRegex.MatchString(item.ShortDescription) {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid short description in item at index %d", i))
		}
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}
	return nil
}

func IsValidPurchaseDate(purchaseDate string) bool {
	_, err := time.Parse("2006-01-02", purchaseDate)
	return err == nil
}
