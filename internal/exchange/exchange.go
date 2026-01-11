package exchange

// Fetcher interface defines methods that any exchange implementation must have.
// like base class in oop
// polymorphism here, the main app doesnt care if its binance or coinbase
type Fetcher interface {
	Name() string
	GetPrice(symbol string) (float64, error)
}

// acts as data transfer object, carries data accross channel
type QuoteResult struct {
	Exchange string
	Price    float64
	Error    error
}
