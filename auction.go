// Package auction implements a computerized auction system where sellers can
// offer items for sale and buyers can place competing bids. The package
// provides mechanisms to start an auction, place bids, and determine the
// winner based on the highest bid. It ensures that bids do not exceed
// maximum limits set by bidders and that each bid is incremented properly
// according to predefined rules.
package auction

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Storer defines the interface for auction data storage operations.
type Storer interface {
	AddBidder(bidder *bidder) error
	GetBidder(id uuid.UUID) (bidder, error)
	UpdateBidder(bidder *bidder) error
	ListBidders() ([]bidder, error)
}

// Auction holds all the details of a single auction event.
type Auction struct {
	storer Storer
	ID     uuid.UUID
}

// NewAuctionConfig is used to configure a new auction.
type NewAuctionConfig struct {
	Bidders []NewBidder
}

// NewAuction creates a new auction instance from the given parameters.
func NewAuction(na NewAuctionConfig) (*Auction, error) {
	if err := validateAuctionData(na); err != nil {
		return nil, fmt.Errorf("invalid auction data: %w", err)
	}

	// -----------------------------------------------------------------------
	// Create bidders and add them to the auction.

	storer := NewInMemoryStore()
	bdrs := toNewBidders(na.Bidders)
	for _, bdr := range bdrs {
		if err := storer.AddBidder(bdr); err != nil {
			return nil, fmt.Errorf("failed to add bidder: %w", err)
		}
	}

	// -----------------------------------------------------------------------

	auction := Auction{
		ID:     uuid.New(),
		storer: storer,
	}

	return &auction, nil
}

// PlaceBid places a bid on the auction.
func (a *Auction) PlaceBid(id uuid.UUID) error {
	bidder, err := a.storer.GetBidder(id)
	if err != nil {
		return fmt.Errorf("failed to get bidder: %w", err)
	}

	bidAmount := bidder.CurrentBid + bidder.AutoIncrement

	// -----------------------------------------------------------------------
	// Perform validations.

	if bidAmount < bidder.StartingBid {
		return fmt.Errorf("bid amount $%.2f is less than starting bid $%.2f", bidAmount, bidder.StartingBid)
	}
	if bidAmount > bidder.MaxBid {
		return fmt.Errorf("%w: bid amount $%.2f is greater than max bid $%.2f", ErrExceededMaxBid, bidAmount, bidder.MaxBid)
	}
	if bidAmount <= bidder.CurrentBid {
		return fmt.Errorf("bid amount $%.2f is less than or equal to current bid $%.2f", bidAmount, bidder.CurrentBid)
	}

	// -----------------------------------------------------------------------
	// Updates the bidder current bid.

	bidder.CurrentBid = bidAmount
	bidder.LastBidTime = time.Now()
	err = a.storer.UpdateBidder(&bidder)
	if err != nil {
		return fmt.Errorf("failed to update bidder: %w", err)
	}

	return nil
}

// Winner represents the winner of the auction.
type Winner struct {
	ID   uuid.UUID
	Name string
}

// DetermineWinner determines the winner of the auction based on the highest current bid.
// In case of a tie (multiple bidders with the same highest bid), the bidder who placed
// their bid first (based on LastBidTime) is considered the winner.
func (a *Auction) DetermineWinner() (Winner, error) {
	bdrs, err := a.storer.ListBidders()
	if err != nil {
		return Winner{}, fmt.Errorf("failed to list bidders: %w", err)
	}

	var wbdr bidder

	for _, bidder := range bdrs {
		isWinner := isWinner(&wbdr, &bidder)
		if isWinner {
			wbdr = bidder
		}
	}

	if wbdr.Name == "" {
		return Winner{}, errors.New("no winner")
	}

	winner := Winner{
		ID:   wbdr.ID,
		Name: wbdr.Name,
	}

	return winner, nil
}

// isWinner checks if the provided bidder should replace the current winner.
// A bidder becomes the new winner if:
// - There is no current winner.
// - Their bid is lower than the current winner's bid.
// - Their bid is the same as the current winner's but was placed earlier.
func isWinner(currentWinner, bidder *bidder) bool {
	return currentWinner.Name == "" || // No current winner, so the bidder wins by default.
		bidder.CurrentBid < currentWinner.CurrentBid || // Bidder has a lower bid.
		(bidder.CurrentBid == currentWinner.CurrentBid && // Bidder has the same bid but placed it earlier.
			bidder.LastBidTime.Before(currentWinner.LastBidTime))
}

// validateAuctionData checks that the provided data for a new auction is valid.
func validateAuctionData(na NewAuctionConfig) error {
	if len(na.Bidders) <= 1 {
		return errors.New("auction must have at least two bidders")
	}

	seenIDs := make(map[uuid.UUID]bool)
	for _, bidder := range na.Bidders {
		// -----------------------------------------------------------------------
		// Check for unique IDs to prevent duplicate bidders.

		if _, exists := seenIDs[bidder.ID]; exists {
			return fmt.Errorf("duplicate bidder ID detected: %s", bidder.ID)
		}
		seenIDs[bidder.ID] = true

		// -----------------------------------------------------------------------
		// Validate individual bidder data.

		if err := validateBidder(&bidder); err != nil {
			return fmt.Errorf("invalid bidder data for bidder ID %s: %w", bidder.ID, err)
		}
	}
	return nil
}

// validateBidder checks that a bidder's data is valid.
func validateBidder(b *NewBidder) error {
	if b.StartingBid <= 0 {
		return fmt.Errorf("starting bid must be positive, got $%.2f", b.StartingBid)
	}
	if b.MaxBid < b.StartingBid {
		return fmt.Errorf("max bid $%.2f must be greater than or equal to starting bid $%.2f", b.MaxBid, b.StartingBid)
	}
	if b.AutoIncrement <= 0 {
		return fmt.Errorf("auto-increment must be positive, got $%.2f", b.AutoIncrement)
	}
	return nil
}
