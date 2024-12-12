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
		validationErrors = append(validationErrors, "The receipt is invalid, retailer is required.")
	}
	if receipt.PurchaseDate == "" {
		validationErrors = append(validationErrors, "The receipt is invalid, purchase date is required.")
	}
	if receipt.PurchaseTime == "" {
		validationErrors = append(validationErrors, "The receipt is invalid, purchase time is required.")
	}
	if receipt.Total == "" {
		validationErrors = append(validationErrors, "The receipt is invalid, total is required.")
	}
	if len(receipt.Items) == 0 {
		validationErrors = append(validationErrors, "The receipt is invalid, item(s) are required.")
	}

	if receipt.Retailer != "" && !retailerRegex.MatchString(receipt.Retailer) {
		validationErrors = append(validationErrors, "The receipt is invalid, bad retailer name.")
	}
	if receipt.PurchaseDate != "" && !IsValidPurchaseDate(receipt.PurchaseDate) {
		validationErrors = append(validationErrors, "The receipt is invalid, bad purchase date. Must be in YYYY-MM-DD format.")
	}
	if receipt.PurchaseTime != "" && !purchaseTimeRegex.MatchString(receipt.PurchaseTime) {
		validationErrors = append(validationErrors, "The receipt is invalid, bad purchase time. Must be in HH:MM format")
	}
	if receipt.Total != "" && !amountRegex.MatchString(receipt.Total) {
		validationErrors = append(validationErrors, "The receipt is invalid, bad total. Must be in ##.## format.")
	}

	for i, item := range receipt.Items {
		if item.Price == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("The receipt is invalid, missing price in item at index %d.", i))
		}
		if item.ShortDescription == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("The receipt is invalid, missing short description in item at index %d.", i))
		}
		if item.Price != "" && !amountRegex.MatchString(item.Price) {
			validationErrors = append(validationErrors, fmt.Sprintf("The receipt is invalid, bad price in item at index %d.", i))
		}
		if item.ShortDescription != "" && !itemShortDescriptionRegex.MatchString(item.ShortDescription) {
			validationErrors = append(validationErrors, fmt.Sprintf("The receipt is invalid, bad short description in item at index %d.", i))
		}
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, " | "))
	}
	return nil
}

func IsValidPurchaseDate(purchaseDate string) bool {
	_, err := time.Parse("2006-01-02", purchaseDate)
	return err == nil
}
