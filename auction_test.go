package auction

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// createBidder helps to create a Bidder with less verbosity.
func createBidder(name string, startingBid, maxBid, increment float64) NewBidder {
	return NewBidder{
		ID:            uuid.New(),
		Name:          name,
		StartingBid:   startingBid,
		MaxBid:        maxBid,
		CurrentBid:    startingBid,
		AutoIncrement: increment,
		LastBidTime:   time.Time{},
	}
}

// TestAuctionScenarios tests multiple auction scenarios.
func TestAuctionScenarios(t *testing.T) {
	tests := []struct {
		name         string
		bidders      []NewBidder
		expectedName string
	}{
		{
			name: "Auction #1",
			bidders: []NewBidder{
				createBidder("Sasha", 50.00, 80.00, 3.00),
				createBidder("John", 60.00, 82.00, 2.00),
				createBidder("Pat", 55.00, 85.00, 5.00),
			},
			expectedName: "Pat",
		},
		{
			name: "Auction #2",
			bidders: []NewBidder{
				createBidder("Riley", 700.00, 725.00, 2.00),
				createBidder("Morgan", 599.00, 725.00, 15.00),
				createBidder("Charlie", 625.00, 725.00, 8.00),
			},
			expectedName: "Riley",
		},
		{
			name: "Auction #3",
			bidders: []NewBidder{
				createBidder("Alex", 2500.00, 3000.00, 500.00),
				createBidder("Jesse", 2800.00, 3100.00, 201.00),
				createBidder("Drew", 2501.00, 3200.00, 247.00),
			},
			expectedName: "Jesse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auction, err := NewAuction(NewAuctionConfig{Bidders: tt.bidders})
			assert.NoError(t, err)

			// Simulate multiple rounds of bidding until no more bids can be placed.
			active := true
			for active {
				active = false
				for _, bidder := range tt.bidders {
					err := auction.PlaceBid(bidder.ID)
					if err != nil && errors.Is(err, ErrExceededMaxBid) {
						break
					}
					if assert.NoError(t, err) {
						active = true // Continue another round if at least one bid was successfully placed.
					}
				}
			}

			winner := auction.DetermineWinner()
			assert.NotNil(t, winner)
			assert.Equal(t, tt.expectedName, winner.Name, "the expected winner does not match.")
		})
	}
}

// TestConcurrencyInAuction tests the PlaceBid function for concurrency issues and data races.
func TestConcurrencyInAuction(t *testing.T) {
	// Setup auction and bidders
	newBidders := []NewBidder{
		createBidder("Sasha", 50.00, 80.00, 3.00),
		createBidder("John", 60.00, 82.00, 2.00),
		createBidder("Pat", 55.00, 85.00, 5.00),
	}

	auction, err := NewAuction(NewAuctionConfig{Bidders: newBidders})
	assert.NoError(t, err)

	// Number of goroutines to simulate concurrent bidding
	numBidders := len(newBidders)
	numBidsPerBidder := 5

	var wg sync.WaitGroup
	wg.Add(numBidders)

	// Launch multiple goroutines to simulate concurrent bidding
	for _, bidder := range newBidders {
		currentBidder := bidder

		go func(nbdr NewBidder) {
			defer wg.Done()

			for i := 0; i < numBidsPerBidder; i++ {
				time.Sleep(time.Millisecond * 10) // simulate delay

				err := auction.PlaceBid(nbdr.ID)
				if err != nil && errors.Is(err, ErrExceededMaxBid) {
					break
				}
				assert.NoError(t, err)
			}
		}(currentBidder)
	}

	wg.Wait()

	// Check if the final bids are within the expected limits
	for _, bidder := range newBidders {
		assert.LessOrEqual(t, bidder.CurrentBid, bidder.MaxBid, "Bid exceeded MaxBid")
	}

	// Determine the winner and ensure it's a valid winner
	winner := auction.DetermineWinner()
	assert.NotNil(t, winner, "There should be a winner")
	assert.NotEmpty(t, winner.Name, "Winner should have a name")
	assert.Equal(t, "Pat", winner.Name, "the expected winner does not match.")
}
