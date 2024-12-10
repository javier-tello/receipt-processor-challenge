package services

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
	"github.com/javier-tello/receipt-processor-challenge/internal/repositories"
)

type ReceiptService struct {
	repo repositories.ReceiptRepository
}

func NewReceiptService(repo repositories.ReceiptRepository) *ReceiptService {
	return &ReceiptService{repo: repo}
}

func (rs *ReceiptService) ProcessReceipt(receipt models.Receipt) (int, error) {
	return rs.repo.ProcessReceipt(receipt), nil
}

func (rs *ReceiptService) GetReceiptByID(id int) (map[string]string, error) {
	receipt, exists := rs.repo.FindByID(id)

	if !exists {
		log.Printf("receipt not found")
		return nil, errors.New("receipt not found")
	}
	points := calculatePointsForDayOfPurchase(receipt) + calculatePointsForItemDescription(receipt) + calculatePointsForItemPairs(receipt) + calculatePointsForRetailerName(receipt) + calculatePointsForTimeOfPurchase(receipt) + calculatePointsForTotalDecimals(receipt)
	data := map[string]string{
		"points": strconv.Itoa(points),
	}

	return data, nil
}

func calculatePointsForRetailerName(receipt models.Receipt) int {
	retailer := receipt.Retailer
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	trimmed_retailer := re.ReplaceAllString(retailer, "")

	return len(trimmed_retailer)
}

func calculatePointsForItemPairs(receipt models.Receipt) int {
	receiptItems := receipt.Items

	return (5 * (len(receiptItems) / 2))
}

func calculatePointsForItemDescription(receipt models.Receipt) int {
	receiptItems := receipt.Items
	points := 0

	for _, value := range receiptItems {
		trimmedDescription := strings.TrimSpace(value.ShortDescription)
		if len(trimmedDescription)%2 == 1 {
			points += 1
		}
		points += (len(trimmedDescription) / 2)
	}

	return points
}

func calculatePointsForTotalDecimals(receipt models.Receipt) int {
	purhcaseTotal := receipt.Total
	points := 0

	indexOfDecimalPoint := strings.IndexRune(purhcaseTotal, '.')
	if purhcaseTotal[indexOfDecimalPoint:] == ".00" || purhcaseTotal[indexOfDecimalPoint:] == ".25" || purhcaseTotal[indexOfDecimalPoint:] == ".75" {
		if purhcaseTotal[indexOfDecimalPoint:] == ".00" {
			points += 50
		}
		points += 25
	}

	return points
}

func calculatePointsForDayOfPurchase(receipt models.Receipt) int {
	purchaseDate := receipt.PurchaseDate
	dayOfPurchase := purchaseDate[8:]
	points := 0

	convertedDayInt, err := strconv.Atoi(dayOfPurchase)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if convertedDayInt%2 == 1 {
			points += 6
		}
	}

	return points
}

func calculatePointsForTimeOfPurchase(receipt models.Receipt) int {
	purchaseTime := receipt.PurchaseTime
	points := 0

	if purchaseTime > "14:00" || purchaseTime < "18:00" {
		points += 10
	}

	return points

}
