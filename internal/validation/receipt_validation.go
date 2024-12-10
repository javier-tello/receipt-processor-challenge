package validation

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
)

type ReceiptValidator struct{}

func (uv *ReceiptValidator) Validate(receipt models.Receipt) error {
	if receipt.Retailer == "" {
		return errors.New("retailer is required")
	}
	if receipt.PurchaseDate == "" {
		return errors.New("purchase date is required")
	}
	if receipt.PurchaseTime == "" {
		return errors.New("purchase time is required")
	}
	if receipt.Total == "" {
		return errors.New("total is required")
	}
	if len(receipt.Items) == 0 {
		return errors.New("item(s) are required")
	}

	if !isValidRetailerName(receipt.Retailer) {
		return errors.New("invalid retailer name, must be alphanumeric characters and the following special characters are allow: & and -")
	}
	if !isValidPurchaseDate(receipt.PurchaseDate) {
		return errors.New("invalid purchase date, must be in YYYY-MM-DD format")
	}
	if !isValidPurchaseTime(receipt.PurchaseTime) {
		return errors.New("invalid time, must be in HH:MM and in military format")
	}
	if !isValidAmount(receipt.Total) {
		return errors.New("invalid total, must be in ##.## format")
	}

	for i := 0; i < len(receipt.Items); i++ {
		if receipt.Items[i].Price == "" {
			return errors.New("price is required and missing in index " + strconv.Itoa(i) + " of items")
		}
		if receipt.Items[i].ShortDescription == "" {
			return errors.New("short description is required and missing in index " + strconv.Itoa(i) + " of items")
		}

		if !isValidAmount(receipt.Items[i].Price) {
			return errors.New("invalid total in index " + strconv.Itoa(i) + " of items, must be in ##.## format")
		}
		if !isValidShortDescription(receipt.Items[i].ShortDescription) {
			return errors.New("invalid short description in index " + strconv.Itoa(i) + " of items, must be alphanumeric characters with - being the only allowed special character")
		}
	}

	return nil
}

func isValidRetailerName(retailer string) bool {
	retailerRegex := `^[\w\s\-&]+$`
	re := regexp.MustCompile((retailerRegex))

	return re.MatchString(retailer)
}

func isValidPurchaseDate(purchaseDate string) bool {
	purchaseDateRegex := `^((\d{4})-(0[13578]|1[02])-(0[1-9]|[12]\d|3[01])|(\d{4})-(0[469]|11)-(0[1-9]|[12]\d|30)|(\d{4})-02-(0[1-9]|1\d|2[0-8])|([02468][048]00|[13579][26]00|[0-9]{2}(0[48]|[2468][048]|[13579][26]))-02-29)$`
	re := regexp.MustCompile(purchaseDateRegex)

	return re.MatchString(purchaseDate)
}

func isValidPurchaseTime(purchaseTime string) bool {
	purchaseTimeRegex := `^(2[0-3]|[01][0-9]):[0-5][0-9]$`
	re := regexp.MustCompile(purchaseTimeRegex)

	return re.MatchString(purchaseTime)
}

func isValidAmount(amount string) bool {
	amountRegx := `^\d+\.\d{2}$`
	re := regexp.MustCompile(amountRegx)

	return re.MatchString(amount)
}

func isValidShortDescription(shortDescription string) bool {
	itemShortDescriptionRegex := `^[\w\s\-]+$`
	re := regexp.MustCompile(itemShortDescriptionRegex)

	return re.MatchString(shortDescription)
}
