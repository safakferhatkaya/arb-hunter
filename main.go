package main

import (
	"fmt"
	"log"
)

type Exchane interface {
	Name() string
	GetPrice(symbol string) (float64, error)
}

type Binance struct{}

func (b Binance) Name() string { return "Binance" }
func (b Binance) GetPrice(symbol string) (float64, error) {
	// Todo: Implement actual API call
	return 98.000, nil
}

type Coinbase struct{}

func (c Coinbase) Name() string { return "Coinbase" }
func (c Coinbase) GetPrice(symbol string) (float64, error) {
	// Todo: Implement actual API call
	return 98.500, nil
}

func main() {
	exchanges := []Exchane{
		Binance{},
		Coinbase{},
	}

	symbol := "BTC"
	fmt.Printf("Searching prices for %s:\n", symbol)

	// Todo : Implement concurrency for fetching prices
	for _, ex := range exchanges {
		price, err := ex.GetPrice(symbol)
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		fmt.Printf("[%s] Price $%.2f\n", ex.Name(), price)
	}
}
