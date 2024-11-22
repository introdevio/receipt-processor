package store

import (
	"errors"
	"github.com/google/uuid"
	"github.com/introdevio/receipt_processor/models"
	"sync"
)

var ErrNotFound = errors.New("receipt not found")

type ReceiptStore interface {
	Save(receipt *models.Receipt) string
	Retrieve(receiptId string) (*models.Receipt, error)
}

type InMemoryReceiptStore struct {
	store map[string]*models.Receipt
	mu    sync.Mutex
}

func NewInMemoryReceiptStore() ReceiptStore {
	return &InMemoryReceiptStore{
		store: make(map[string]*models.Receipt),
		mu:    sync.Mutex{},
	}
}

func (s *InMemoryReceiptStore) Save(receipt *models.Receipt) string {
	id := uuid.NewString()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[id] = receipt
	return id
}

func (s *InMemoryReceiptStore) Retrieve(id string) (*models.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if receipt, exists := s.store[id]; exists {
		return receipt, nil
	} else {
		return nil, ErrNotFound
	}
}
