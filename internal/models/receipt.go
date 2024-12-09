package models

// Contents of a receipt
type Receipt struct {
	ID           int
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Repository
type ReceiptRepository interface {
	FindByID(id int) (*Receipt, error)
	Save(receipt *Receipt) error
}

// In memeory implemetaion for challenge
type InMemoryReceiptRepo struct {
	receipts map[int]*Receipt
	nextID   int
}

func NewInMemoryReceiptRepo() *InMemoryReceiptRepo {
	return &InMemoryReceiptRepo{
		receipts: make(map[int]*Receipt),
		nextID:   1,
	}
}

// GET and POST calls
func (repo *InMemoryReceiptRepo) FindByID(id int) (*Receipt, error) {
	receipt, exists := repo.receipts[id]
	if !exists {
		return nil, nil
	}
	return receipt, nil
}

func (repo *InMemoryReceiptRepo) Save(receipt *Receipt) error {
	if receipt.ID == 0 {
		receipt.ID = repo.nextID
		repo.nextID++
	}
	repo.receipts[receipt.ID] = receipt
	return nil
}
