package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type coinbase struct{}

// Factory method
func NewCoinbase() Fetcher {
	return &coinbase{}
}

func (c *coinbase) Name() string { return "Coinbase" }

// coinbase api response structure is nested
// {"data":{"base":"BTC", "currency":"USD", "amount":"9800.00"}}
type coinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

func (c *coinbase) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-USD/spot", symbol)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data coinbaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(data.Data.Amount, 64)
	if err != nil {
		return 0, fmt.Errorf("Can't convert price: %v", err)
	}

	return price, nil
}
