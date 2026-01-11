package analyzer

import (
	"arb-hunter/internal/exchange"
	"testing"
)

func TestCalculate(t *testing.T) {
	// Table-Driven Test: The Go standard for testing logic.
	// We define a list of scenarios (structs) and loop through them.
	// In RSpec, these would be separate 'context' blocks.
	tests := []struct {
		name           string                 // Name of the test case
		input          []exchange.QuoteResult // Input data
		expectedSpread float64                // Expected result
		expectError    bool                   // Do we expect an error?
	}{
		{
			name: "Not enough data (Single Exchange)",
			input: []exchange.QuoteResult{
				{Exchange: "Binance", Price: 50000},
			},
			expectError: true, // We need at least 2 to compare
		},
		{
			name: "Profitable Arbitrage (Buy Binance, Sell Coinbase)",
			input: []exchange.QuoteResult{
				{Exchange: "Binance", Price: 100},  // Min
				{Exchange: "Coinbase", Price: 110}, // Max
			},
			// Math: (110 - 100) / 100 * 100 = 10%
			expectedSpread: 10.0,
			expectError:    false,
		},
		{
			name: "No Profit (Prices are equal)",
			input: []exchange.QuoteResult{
				{Exchange: "Binance", Price: 100},
				{Exchange: "Coinbase", Price: 100},
			},
			expectedSpread: 0.0,
			expectError:    false,
		},
		{
			name: "Negative Spread (Market Inefficiency)",
			// Even if the "first" one is expensive, our logic sorts finding Min/Max.
			// So it will effectively treat this as: Buy Coinbase (100), Sell Binance (120).
			input: []exchange.QuoteResult{
				{Exchange: "Binance", Price: 120},
				{Exchange: "Coinbase", Price: 100},
			},
			expectedSpread: 20.0,
			expectError:    false,
		},
	}

	// The Runner Loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := Calculate(tt.input)

			// 1. Check Error State
			if (err != nil) != tt.expectError {
				t.Errorf("Calculate() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// If we expected an error and got one, stop here.
			if tt.expectError {
				return
			}

			// 2. Check Math (Spread)
			if report.Spread != tt.expectedSpread {
				t.Errorf("Calculate() spread = %v, want %v", report.Spread, tt.expectedSpread)
			}
		})
	}
}
