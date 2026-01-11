package analyzer

import (
	"arb-hunter/internal/exchange"
	"fmt"
)

type Report struct {
	BestBuy  exchange.QuoteResult
	BestSell exchange.QuoteResult
	Spread   float64
}

// Calculate contains the core business logic

func Calculate(quotes []exchange.QuoteResult) (Report, error) {
	if len(quotes) < 2 {
		return Report{}, fmt.Errorf("not enough data to perform analsis")
	}

	min := quotes[0]
	max := quotes[0]

	for _, quote := range quotes[1:] {
		if quote.Price < min.Price {
			min = quote
		}
		if quote.Price > max.Price {
			max = quote
		}
	}

	spread := (max.Price - min.Price) / min.Price * 100

	return Report{
		BestBuy:  min,
		BestSell: max,
		Spread:   spread,
	}, nil
}
