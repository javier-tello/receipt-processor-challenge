package services

import (
	"fmt"
	"log"
	"math"
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

func (rs *ReceiptService) CalculateTotalPointsForReceipt(receiptID int) (int, error) {
	log.Println("Retreiving receipt")
	receipt, exists := rs.repo.FindByID(receiptID)
	if !exists {
		log.Println("receipt not found")
	}
	log.Println("Receipt successfully retreived")

	log.Println("Calculating points for receipt")
	points := calculatePointsForDayOfPurchase(receipt.PurchaseDate) + calculatePointsForItemDescription(receipt.Items) + calculatePointsForItemPairs(receipt.Items) + calculatePointsForRetailerName(receipt.Retailer) + calculatePointsForTimeOfPurchase(receipt.PurchaseTime) + calculatePointsForTotalDecimals(receipt.Total)

	return points, nil
}

func calculatePointsForRetailerName(retailerName string) int {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	trimmed_retailer := re.ReplaceAllString(retailerName, "")

	log.Println(len(trimmed_retailer), " points - retailer name (", trimmed_retailer, ") has ", len(trimmed_retailer), " alphanumeric")

	return len(trimmed_retailer)
}

func calculatePointsForItemPairs(items []models.Item) int {
	log.Println(5*(len(items)/2), " points - ", len(items), " items (2 pairs @ 5 points each)")

	return (5 * (len(items) / 2))
}

func calculatePointsForItemDescription(items []models.Item) int {
	points := 0

	for _, value := range items {
		trimmedDescription := strings.TrimSpace(value.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			price, err := strconv.ParseFloat(value.Price, 64)
			if err != nil {
				log.Println("Error in converting string to int")
			}
			points += int(math.Ceil(price * 0.2))
			log.Println(int(math.Ceil(price*0.2)), " Points - \"", trimmedDescription, "\" is ", len(trimmedDescription), " characters (a multiple of 3) item price of ", value.Price, " * 0.2 = ", price*0.2, ", rounded up is ", int(math.Ceil(price*0.2)), " points")
		}
	}

	return points
}

func calculatePointsForTotalDecimals(purhcaseTotal string) int {
	points := 0
	indexOfDecimalPoint := strings.IndexRune(purhcaseTotal, '.')
	if purhcaseTotal[indexOfDecimalPoint:] == ".00" || purhcaseTotal[indexOfDecimalPoint:] == ".25" || purhcaseTotal[indexOfDecimalPoint:] == ".75" {
		if purhcaseTotal[indexOfDecimalPoint:] == ".00" {
			points += 50
			log.Println("50 points - total is a round dollar amount")
		}
		points += 25
		log.Println("25 points - total is a multiple of 0.25")
	}

	return points
}

func calculatePointsForDayOfPurchase(purchaseDate string) int {
	dayOfPurchase := purchaseDate[8:]
	points := 0

	convertedDayInt, err := strconv.Atoi(dayOfPurchase)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		if convertedDayInt%2 == 1 {
			points += 6
			log.Println("6 points - purchase day is odd")
		}
	}

	return points
}

func calculatePointsForTimeOfPurchase(purchaseTime string) int {
	points := 0

	if purchaseTime > "14:00" && purchaseTime < "18:00" {
		points += 10
		log.Println("10 points - ", purchaseTime, " is between 14:00 and 16:00")
	}

	return points

}
