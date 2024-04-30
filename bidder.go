package auction

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrExceededMaxBid is returned when a bid exceeds the maximum bid allowed by a bidder.
var ErrExceededMaxBid = errors.New("exceeded maximum bid")

// Bidder represents an individual participant in an auction.
type bidder struct {
	ID            uuid.UUID
	Name          string
	StartingBid   float64
	MaxBid        float64
	CurrentBid    float64
	AutoIncrement float64
	LastBidTime   time.Time
}

// NewBidder is used to configure a new bidder.
type NewBidder struct {
	ID            uuid.UUID
	Name          string
	StartingBid   float64
	MaxBid        float64
	CurrentBid    float64
	AutoIncrement float64
	LastBidTime   time.Time
}

// toNewBidders converts a slice of NewBidder to a slice of bidder.
func toNewBidders(bidders []NewBidder) []*bidder {
	var newBidders []*bidder
	for _, bdr := range bidders {
		newBidder := toNewBidder(bdr)
		newBidders = append(newBidders, &newBidder)
	}
	return newBidders
}

// toNewBidder converts a NewBidder to a bidder.
func toNewBidder(bdr NewBidder) bidder {
	return bidder{
		ID:            bdr.ID,
		Name:          bdr.Name,
		StartingBid:   bdr.StartingBid,
		MaxBid:        bdr.MaxBid,
		CurrentBid:    bdr.CurrentBid,
		AutoIncrement: bdr.AutoIncrement,
		LastBidTime:   bdr.LastBidTime,
	}
}
