# 💃 Lady Whistle - Stock Predictor Agent 📈

A lightweight CLI tool written in **Go** that fetches real stock market data from Yahoo Finance, trains a simple linear regression model, and predicts the closing price for the next 7 trading days.

> **Disclaimer:** This is an educational project. The linear regression model is intentionally simple and does **not** account for market sentiment, news, earnings, or non-linear dynamics. Do not use it for real trading decisions.

---

## Features

- **Live data** — pulls up to 6 months of daily OHLCV data directly from Yahoo Finance
- **Linear regression** — trains an OLS (Ordinary Least Squares) model on closing prices
- **7-day forecast** — predicts the next 7 trading days (weekends are skipped)
- **Market summary** — displays SMA-20, SMA-50, volatility, period high/low
- **Model quality** — reports R² and warns when the trend is weak
- **Confidence band** — shows ±1σ range around the final prediction
- **Zero dependencies** — only uses the Go standard library

---

## Project Structure

```
Lady-Whistle/
├── main.go        # Entry point – CLI args, orchestration, output
├── yahoo.go       # Yahoo Finance API types and data fetching
├── model.go       # Linear regression training and prediction
├── analysis.go    # Statistical helpers (SMA, volatility, min/max)
├── display.go     # Terminal formatting (headers, dividers)
├── go.mod         # Go module definition
└── README.md      # This file
```

| File            | Purpose                                      |
| --------------- | -------------------------------------------- |
| `main.go`       | Parses CLI args, orchestrates fetch → train → predict, prints results |
| `yahoo.go`      | Defines API response structs and `FetchStockData()` |
| `model.go`      | `LinearModel` struct, `Train()`, and `Predict()` |
| `analysis.go`   | `SimpleMovingAverage()`, `Volatility()`, `MaxPrice()`, `MinPrice()` |
| `display.go`    | `PrintHeader()`, `PrintDivider()` for clean terminal output |

---

## Quick Start

### Prerequisites

- Internet connection (for Yahoo Finance API)

### Build & Run

```bash

# Build the binary
go build -o predict

# Run with default ticker (AAPL)
./predict

# Run with a specific ticker
./predict TSLA
./predict MSFT
./predict GOOG
./predict AMZN
```

### Example Output

```
════════════════════════════════════════════════════════════
  Lady Whistle   ·  TSLA
════════════════════════════════════════════════════════════
  Fetching 180 days of historical data …
  ✓ Received 124 trading days of data

════════════════════════════════════════════════════════════
  Market Summary
════════════════════════════════════════════════════════════
  Latest close : $405.50  (05-03-2026)
  Period high  : $489.88
  Period low   : $346.40
  SMA-20       : $410.20
  SMA-50       : $429.85
  Volatility   : $8.67  (20-day std dev)

════════════════════════════════════════════════════════════
  Model Training
════════════════════════════════════════════════════════════
  Model        : Linear Regression (OLS)
  Slope        : 0.0062  ($/day)
  Intercept    : $431.27
  R²           : 0.0001
  ⚠  Low R² – the trend is weak; predictions may be unreliable.

════════════════════════════════════════════════════════════
  7-Day Price Forecast for TSLA
════════════════════════════════════════════════════════════
  Date          Predicted     Change
────────────────────────────────────────────────────────────
  06-03-2026    $432.04       ▲ +26.54 (+6.55%)
  09-03-2026    $432.05       ▲ +26.55 (+6.55%)
  ...
────────────────────────────────────────────────────────────

  Confidence band (±1σ): $423.41 – $440.75

  ⚠  Disclaimer: This is a simple linear regression model.
```

---

## How It Works

1. **Fetch** — `yahoo.go` sends an HTTP GET request to Yahoo Finance's chart API and parses the JSON response into `[]StockData`
2. **Analyze** — `analysis.go` computes indicators like SMA and volatility from the closing prices
3. **Train** — `model.go` fits a linear regression (y = mx + b) using ordinary least squares on the historical closing prices
4. **Predict** — The trained model extrapolates the trend line for the next 7 trading days
5. **Display** — `main.go` and `display.go` format everything into a clean terminal report

---

## Limitations

- **Linear model only** — cannot capture non-linear trends, mean reversion, or momentum
- **No feature engineering** — uses raw day-index as the only feature
- **No volume/sentiment data** — ignores trading volume, news, and other signals
- **Weekend-only skip** — does not account for market holidays
- **Past performance ≠ future results** — this is a fundamental limitation of all models

---

## License

MIT — use freely for learning and experimentation.
