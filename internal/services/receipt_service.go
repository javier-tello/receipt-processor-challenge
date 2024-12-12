package services

import (
	"errors"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
	"github.com/javier-tello/receipt-processor-challenge/internal/repositories"
)

const (
	PointsForOddDay          = 6
	PointsForTimeWindow      = 10
	PointsForRoundDollar     = 50
	PointsForQuarterMultiple = 25
	PointsPerItemPair        = 5
)

type ReceiptService struct {
	repo repositories.ReceiptRepository
}

func NewReceiptService(repo repositories.ReceiptRepository) *ReceiptService {
	return &ReceiptService{repo: repo}
}

func (rs *ReceiptService) ProcessReceipt(receipt models.Receipt) string {
	return rs.repo.ProcessReceipt(receipt)
}

func (rs *ReceiptService) CalculateTotalPointsForReceipt(receiptID string) (int, error) {
	log.Println("Retreiving receipt")
	receipt, exists := rs.repo.FindByID(receiptID)
	if !exists {
		return -1, errors.New("cannot find receipt")
	}
	log.Println("Receipt successfully retreived")

	log.Println("Calculating points for receipt")
	points := calculatePointsForDayOfPurchase(receipt.PurchaseDate) +
		calculatePointsForItemDescription(receipt.Items) +
		calculatePointsForItemPairs(receipt.Items) +
		calculatePointsForRetailerName(receipt.Retailer) +
		calculatePointsForTimeOfPurchase(receipt.PurchaseTime) +
		calculatePointsForTotalDecimals(receipt.Total)

	return points, nil
}

func calculatePointsForRetailerName(retailerName string) int {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	trimmed_retailer := re.ReplaceAllString(retailerName, "")

	log.Println("\t", len(trimmed_retailer), " points - retailer name (", trimmed_retailer, ") has ", len(trimmed_retailer), " alphanumeric")

	return len(trimmed_retailer)
}

func calculatePointsForItemPairs(items []models.Item) int {
	log.Println("\t", 5*(len(items)/2), " points - ", len(items), " items (2 pairs @ 5 points each)")

	return (PointsPerItemPair * (len(items) / 2))
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
			log.Println("\t", int(math.Ceil(price*0.2)), " Points - \"", trimmedDescription, "\" is ", len(trimmedDescription), " characters (a multiple of 3) item price of ", value.Price, " * 0.2 = ", price*0.2, ", rounded up is ", int(math.Ceil(price*0.2)), " points")
		}
	}

	return points
}

func calculatePointsForTotalDecimals(purchaseTotal string) int {
	total, err := strconv.ParseFloat(purchaseTotal, 64)
	if err != nil {
		log.Printf("Invalid total: %v", err)
		return 0
	}

	points := 0
	if math.Mod(total, 1.0) == 0 {
		points += PointsForRoundDollar
		log.Println("\t50 points - total is a round dollar amount")
	}
	if math.Mod(total, 0.25) == 0 {
		points += PointsForQuarterMultiple
		log.Println("\t25 points - total is a multiple of 0.25")
	}

	return points
}

func calculatePointsForDayOfPurchase(purchaseDate string) int {
	day := purchaseDate[8:]

	if day[len(day)-1]%2 == 1 {
		log.Println("\t6 points - purchase day is odd")
		return PointsForOddDay
	}

	return 0
}

func calculatePointsForTimeOfPurchase(purchaseTime string) int {
	if purchaseTime > "14:00" && purchaseTime < "18:00" {
		log.Println("\t10 points - ", purchaseTime, " is between 14:00 and 16:00")
		return PointsForTimeWindow
	}

	return 0

}
