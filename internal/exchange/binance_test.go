package exchange

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBinance_GetPrice(t *testing.T) {
	// 1. SETUP MOCK SERVER (The "VCR" of Go)
	// This spins up a local server on a random port to simulate Binance API.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Contract Testing: Verify we are sending the correct request
		expectedPath := "/api/v3/ticker/price"
		if r.URL.Path != expectedPath {
			t.Errorf("Wrong path. Expected: %s, Got: %s", expectedPath, r.URL.Path)
		}

		if r.URL.Query().Get("symbol") != "BTCUSDT" {
			t.Errorf("Wrong query param. Expected BTCUSDT, Got: %s", r.URL.Query().Get("symbol"))
		}

		// Mock Response: Success case
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"symbol":"BTCUSDT","price":"98500.50"}`)
	}))

	// Teardown: Close the server when test finishes to free up the port.
	defer server.Close()

	// 2. INITIALIZATION (White-box)
	// Since we are in the 'exchange' package, we can initialize the private struct directly.
	// We inject the 'server.URL' (e.g., http://127.0.0.1:5321) into baseURL.
	b := &binance{
		baseURL: server.URL,
	}

	// 3. EXECUTE
	price, err := b.GetPrice("BTC")

	// 4. ASSERT
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	expectedPrice := 98500.50
	if price != expectedPrice {
		t.Errorf("Price mismatch. Expected: %f, Got: %f", expectedPrice, price)
	}
}

func TestBinance_GetPrice_InvalidJSON(t *testing.T) {
	// Scenario: Server returns garbage data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"symbol": "BTC", "price": "INVALID_NUMBER"}`)
	}))
	defer server.Close()

	// Injecting the mock server URL
	b := &binance{
		baseURL: server.URL,
	}

	_, err := b.GetPrice("BTC")

	// We expect an error here because ParseFloat should fail.
	if err == nil {
		t.Error("Expected JSON parsing error, but got nil")
	}
}
