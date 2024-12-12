package repositories

import (
	"sync"

	"github.com/javier-tello/receipt-processor-challenge/internal/models"

	"github.com/google/uuid"
)

// UUIDGenerator is an interface for generating UUIDs
type UUIDGenerator interface {
	New() uuid.UUID
}

// DefaultUUIDGenerator uses the google/uuid package to generate UUIDs
type DefaultUUIDGenerator struct{}

func (g DefaultUUIDGenerator) New() uuid.UUID {
	return uuid.New()
}

// Repository
type ReceiptRepository interface {
	ProcessReceipt(receipt models.Receipt) string
	FindByID(id string) (models.Receipt, bool)
}

// In-memory implementation for this challenge
type InMemoryReceiptRepo struct {
	receipts    map[string]models.Receipt
	idGenerator UUIDGenerator
	mu          sync.RWMutex
}

func NewInMemoryReceiptRepo(generator UUIDGenerator) *InMemoryReceiptRepo {
	if generator == nil {
		generator = DefaultUUIDGenerator{}
	}
	return &InMemoryReceiptRepo{
		receipts:    make(map[string]models.Receipt),
		idGenerator: generator,
	}
}

// FindByID retrieves a receipt by its ID. Returns the receipt and a boolean indicating if it exists.
func (repo *InMemoryReceiptRepo) FindByID(receiptID string) (models.Receipt, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	receipt, ok := repo.receipts[receiptID]

	return receipt, ok
}

// ProcessReceipt saves a receipt in memory and returns its generated ID.
func (repo *InMemoryReceiptRepo) ProcessReceipt(receipt models.Receipt) string {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	receiptID := repo.idGenerator.New().String()

	repo.receipts[receiptID] = receipt
	return receiptID
}
