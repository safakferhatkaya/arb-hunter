# ArbHunter

**ArbHunter** is a CLI tool built specifically to **explore and learn Go (Golang)** concepts.

It simulates a crypto arbitrage detector to demonstrate how Go handles concurrency, interfaces, and project structure compared to object-oriented languages like Ruby.

## 📚 Learning Outcomes

This project serves as a practical implementation of the following Go concepts:

* **Concurrency vs. Parallelism:** Moving from sequential HTTP requests to concurrent fetching using **Goroutines** and **Channels**.
* **Sync Patterns:** Using `sync.WaitGroup` to orchestrate multiple workers.
* **Standard Project Layout:** Structuring a project with `cmd/` and `internal/` directories instead of a flat file.
* **Interfaces:** Abstraction of external dependencies (Binance/Coinbase) using the `Fetcher` interface.
* **Testing:** Implementing **Table-Driven Tests** for business logic.

## 📂 Project Structure

The project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```text
arb-hunter/
├── cmd/
│   └── arb-hunter/
│       └── main.go       # Entry point (Wiring & Orchestration)
├── internal/
│   ├── analyzer/         # Pure Business Logic (Calculates spread)
│   └── exchange/         # API Implementations & Interfaces
├── go.mod
└── README.md
```
🛠️ Usage
Run the CLI
Fetches prices from Binance and Coinbase concurrently.

```Bash
go mod tidy
go run cmd/arb-hunter/main.go
```

Run Tests
Executes table-driven unit tests for the analyzer logic.

```Bash
go test ./... -v
```

📝 Note
This is a playground project created for educational purposes. It is not intended for production-grade financial trading.****
