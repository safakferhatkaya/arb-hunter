package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Exchange interface {
	Name() string
	GetPrice(symbol string) (float64, error)
}

// Implementation of binance
type Binance struct{}

func (b Binance) Name() string { return "Binance" }

// Struct for binance api - Api: {"symbol":"BTCUSDT","price":"9800.00"}
type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"` // Price as string to match API response
}

func (b Binance) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data BinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("Can't convert price: %v", err)
	}

	return price, nil
}

// Implementation of coinbase
type Coinbase struct{}

func (c Coinbase) Name() string { return "Coinbase" }

// Struct for coinbase api - Api: {"data":{"base":"BTC", "currency":"USD", "amount":"9800.00"}}
type CoinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

func (c Coinbase) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-USD/spot", symbol)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var response CoinbaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	price, err := strconv.ParseFloat(response.Data.Amount, 64)
	if err != nil {
		return 0, fmt.Errorf("Can't convert price: %v", err)
	}
	return price, nil
}

// Struct to hold quote result
type QuoteResult struct {
	Exchange string
	Price    float64
	Error    error
}

func main() {
	start := time.Now()

	defer func() {
		fmt.Printf("Elapsed time: %s\n", time.Since(start))
	}()

	exchanges := []Exchange{
		Binance{},
		Coinbase{},
	}

	symbol := "BTC"
	fmt.Printf("Searching prices for %s:\n", symbol)
	fmt.Println("-----------------------")

	results := make(chan QuoteResult, len(exchanges))
	var wg sync.WaitGroup

	// Start a worker for each market
	for _, ex := range exchanges {
		wg.Add(1)

		go func(e Exchange) {
			defer wg.Done() // Works after current goroutine is done

			price, err := e.GetPrice(symbol)

			// Pass result into channel
			results <- QuoteResult{
				Exchange: e.Name(),
				Price:    price,
				Error:    err,
			}
		}(ex)
	}

	// Closer, take it as seperate goroutine to avoid blocking
	go func() {
		wg.Wait()      // Wait for all fetches to complete
		close(results) // Close the channel
	}()

	var finalQuotes []QuoteResult

	// Consumer, until channel is closed, read results and prints
	for res := range results {
		if res.Error != nil {
			fmt.Printf("[%s] Error fetching price: %v\n", res.Exchange, res.Error)
			continue
		}

		finalQuotes = append(finalQuotes, res)
	}

	if len(finalQuotes) < 2 {
		fmt.Println("Not enough quotes to compare.")
		return
	}

	var minQuote, maxQuote QuoteResult
	minQuote = finalQuotes[0]
	maxQuote = finalQuotes[0]

	for _, quote := range finalQuotes[1:] {
		if quote.Price < minQuote.Price {
			minQuote = quote
		}
		if quote.Price > maxQuote.Price {
			maxQuote = quote
		}
	}

	spread := (maxQuote.Price - minQuote.Price) / minQuote.Price * 100

	fmt.Printf("Lowest Price: [%s] $%.2f\n", minQuote.Exchange, minQuote.Price)
	fmt.Printf("Highest Price: [%s] $%.2f\n", maxQuote.Exchange, maxQuote.Price)
	fmt.Println("\n=== ARBITRAGE REPORT ===")
	fmt.Printf("Best Buy : %s ($%.2f)\n", minQuote.Exchange, minQuote.Price)
	fmt.Printf("Best Sell: %s ($%.2f)\n", maxQuote.Exchange, maxQuote.Price)
	fmt.Printf("Profit   : %.2f%%\n", spread)
	fmt.Println("========================")
}
