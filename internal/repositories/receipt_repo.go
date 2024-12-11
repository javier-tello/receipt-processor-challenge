package repositories

import (
	"sync"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"
)

// Repository
type ReceiptRepository interface {
	ProcessReceipt(receipt models.Receipt) int
	FindByID(id int) (models.Receipt, bool)
}

// In memeory implemetaion for this challenge
type InMemoryReceiptRepo struct {
	receipts  map[int]models.Receipt
	idCounter int
	mu        sync.Mutex
}

func NewInMemoryReceiptRepo() *InMemoryReceiptRepo {
	return &InMemoryReceiptRepo{
		receipts:  make(map[int]models.Receipt),
		idCounter: 0,
	}
}

func (repo *InMemoryReceiptRepo) FindByID(id int) (models.Receipt, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	receipt, exists := repo.receipts[id]
	return receipt, exists
}

func (repo *InMemoryReceiptRepo) ProcessReceipt(receipt models.Receipt) int {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	receipt.ID = repo.idCounter
	repo.idCounter++

	repo.receipts[receipt.ID] = receipt
	return receipt.ID
}
