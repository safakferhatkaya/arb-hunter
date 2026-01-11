package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/safakferhatkaya/arb-hunter/internal/analyzer"
	"github.com/safakferhatkaya/arb-hunter/internal/exchange"
)

func main() {
	start := time.Now()

	// Run at end of the main func
	defer func() {
		fmt.Printf("Elapsed time: %s\n", time.Since(start))
	}()

	// Dependency injection
	exchanges := []exchange.Fetcher{
		exchange.NewBinance(),
		exchange.NewCoinbase(),
	}

	symbol := "BTC"
	fmt.Printf("Searching prices for %s:\n", symbol)
	fmt.Println("-----------------------")

	// Concurrency setup
	results := make(chan exchange.QuoteResult, len(exchanges))
	var wg sync.WaitGroup

	// Start a worker for each market
	for _, ex := range exchanges {
		wg.Add(1)

		go func(e exchange.Fetcher) {
			defer wg.Done() // Works after current goroutine is done

			price, err := e.GetPrice(symbol)

			// Pass result into channel
			results <- exchange.QuoteResult{
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

	// Agggregate results
	var finalQuotes []exchange.QuoteResult
	for res := range results {
		if res.Error != nil {
			log.Printf("[%s] Error: %v\n", res.Exchange, res.Error)
			continue
		}
		fmt.Printf("[%s] Price: $%.2f\n", res.Exchange, res.Price)
		finalQuotes = append(finalQuotes, res)
	}

	// Business logic execution
	report, err := analyzer.Calculate(finalQuotes)
	if err != nil {
		log.Println("Analysis skipped: ", err)
	}

	// Presentation Layer
	fmt.Println("\n=== ARBITRAGE REPORT ===")
	fmt.Printf("Best Buy : %s ($%.2f)\n", report.BestBuy.Exchange, report.BestBuy.Price)
	fmt.Printf("Best Sell: %s ($%.2f)\n", report.BestSell.Exchange, report.BestSell.Price)
	fmt.Printf("Profit   : %.2f%%\n", report.Spread)
	fmt.Println("========================")

	fmt.Println("-----------------------")
}
