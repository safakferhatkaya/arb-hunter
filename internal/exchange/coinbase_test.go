package exchange

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoinbase_GetPrice(t *testing.T) {
	// 1. MOCK SERVER
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate Path (Coinbase API structure is different)
		// Expected: /v2/prices/BTC-USD/spot
		expectedPath := "/v2/prices/BTC-USD/spot"
		if r.URL.Path != expectedPath {
			t.Errorf("Wrong path. Expected: %s, Got: %s", expectedPath, r.URL.Path)
		}

		// Mock Response: Nested JSON structure
		w.WriteHeader(http.StatusOK)
		// Pay attention to the JSON structure matching 'coinbaseResponse' struct
		fmt.Fprintln(w, `{"data": {"base":"BTC", "currency":"USD", "amount":"99100.00"}}`)
	}))
	defer server.Close()

	// 2. INITIALIZE
	// Injecting the Mock Server URL into the private struct
	c := &coinbase{
		baseURL: server.URL,
	}

	// 3. EXECUTE
	price, err := c.GetPrice("BTC")

	// 4. ASSERT
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedPrice := 99100.00
	if price != expectedPrice {
		t.Errorf("Price mismatch. Expected: %f, Got: %f", expectedPrice, price)
	}
}
