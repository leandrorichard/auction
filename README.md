
# Dispatch Bidder Auction System

## Overview
The Dispatch Bidder Auction System is a Go library designed to handle automated bidding processes for a computerized auction platform. This system allows sellers to offer items for sale and enables buyers to compete by placing bids.

## Features
- **Starting Bid:** The initial amount a buyer is willing to offer for an item.
- **Max Bid:** The maximum amount a bidder is prepared to pay for the item.
- **Auto-Increment:** A specified amount added to the current bid when a bidder is outbid by another.
- **Bid Cap:** Ensures that the current bid does not exceed the maximum bid specified by a bidder.
- **Minimum Winning Bid:** Calculates the smallest possible winning bid, adhering to the auto-increment rule.
- **Tie Resolution:** Prioritizes bidders based on who placed the bid first if there is a tie.

## Usage
### Creating an Auction
To start an auction, you need to configure it with bidders, each having a unique identifier, starting bid, maximum bid, and an auto-increment value.

```go
import (
    "github.com/leandrorichard/auction"
    "github.com/google/uuid"
)

func main() {
    bidders := []*auction.Bidder{
        {
            ID: uuid.New(),
            Name: "John Doe",
            StartingBid: 50.0,
            MaxBid: 100.0,
            AutoIncrement: 5.0,
        },
        {
            ID: uuid.New(),
            Name: "Jane Smith",
            StartingBid: 60.0,
            MaxBid: 120.0,
            AutoIncrement: 10.0,
        },
    }
    
    auctionConfig := auction.NewAuctionConfig{Bidders: bidders}
    auction, err := auction.NewAuction(auctionConfig)
    if err != nil {
        fmt.Println("error creating auction:", err)
        return
    }
    
    // Bidding
    err = auction.PlaceBid(bidders[0], 55)
    if err != nil {
        fmt.Println("error placing bid:", err)
    }
    
    // Determining the winner
    winner := auction.DetermineWinner()
    fmt.Printf("winner is %s with a bid of $%.2f", winner.Name, winner.CurrentBid)
}
```

## Installation
To use the Dispatch Bidder Auction System, ensure you have Go installed on your machine. You can then add the library to your project:
```bash
go get github.com/leandrorichard/auction
```

## Testing
Run tests using the Go command:
```bash
make test && make test-race
```