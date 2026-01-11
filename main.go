package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Exchane interface {
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

func main() {
	exchanges := []Exchane{
		Binance{},
		Coinbase{},
	}

	symbol := "BTC"
	fmt.Printf("Searching prices for %s:\n", symbol)
	fmt.Println("-----------------------")

	start := time.Now()

	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("Elapsed time: %s\n", elapsed)
	}()

	// Todo : Implement concurrency for fetching prices
	for _, ex := range exchanges {
		price, err := ex.GetPrice(symbol)
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		fmt.Printf("[%s] Price $%.2f\n", ex.Name(), price)
	}

	fmt.Println("-----------------------")
}
