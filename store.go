package auction

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// InMemoryStore implements the BidderStore interface using an in-memory map.
type InMemoryStore struct {
	sync.RWMutex
	bidders map[uuid.UUID]*bidder
}

// NewInMemoryStore creates a new InMemoryStore instance.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		bidders: make(map[uuid.UUID]*bidder),
	}
}

// AddBidder adds a new bidder to the store.
func (store *InMemoryStore) AddBidder(bidder *bidder) error {
	store.Lock()
	defer store.Unlock()
	if _, exists := store.bidders[bidder.ID]; exists {
		return errors.New("bidder already exists")
	}
	store.bidders[bidder.ID] = bidder
	return nil
}

// GetBidder retrieves a bidder from the store.
func (store *InMemoryStore) GetBidder(id uuid.UUID) (bidder, error) {
	store.RLock()
	defer store.RUnlock()
	bdr, exists := store.bidders[id]
	if !exists {
		return bidder{}, errors.New("bidder not found")
	}
	return *bdr, nil
}

// UpdateBidder updates a bidder in the store.
func (store *InMemoryStore) UpdateBidder(bidder *bidder) error {
	store.Lock()
	defer store.Unlock()
	store.bidders[bidder.ID] = bidder
	return nil
}

// ListBidders retrieves all bidders from the store.
func (store *InMemoryStore) ListBidders() ([]bidder, error) {
	store.RLock()
	defer store.RUnlock()
	bidders := make([]bidder, 0, len(store.bidders))
	for _, bidder := range store.bidders {
		bidders = append(bidders, *bidder)
	}
	return bidders, nil
}
