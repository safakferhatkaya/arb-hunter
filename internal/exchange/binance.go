package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// binance is unexported (private)
type binance struct{}

func NewBinance() Fetcher {
	return &binance{}
}

func (b *binance) Name() string { return "Binance" }

type binanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (b *binance) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data binanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("Can't convert price: %v", err)
	}

	return price, nil
}
