package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	ticker := "AAPL"
	if len(os.Args) > 1 {
		ticker = strings.ToUpper(os.Args[1])
	}
	ticker = strings.TrimSpace(ticker)

	lookbackDays := 180 // ~6 months of history, can be changed
	forecastDays := 7

	// Fetch data
	PrintHeader(fmt.Sprintf("Lady Whistle  ·  %s", ticker))
	fmt.Printf("  Fetching %d days of historical data …\n", lookbackDays)

	data, err := FetchStockData(ticker, lookbackDays)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n  ✗ Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ Received %d trading days of data\n", len(data))

	if len(data) < 10 {
		fmt.Fprintln(os.Stderr, "\n  ✗ Not enough data to train a model.")
		os.Exit(1)
	}

	// Extract closing prices
	prices := make([]float64, len(data))
	for i, d := range data {
		prices[i] = d.Close
	}

	// Market summary 
	latest := data[len(data)-1]
	PrintHeader("Market Summary")
	fmt.Printf("  Latest close : $%.2f  (%s)\n", latest.Close, latest.Date.Format("02-01-2006"))
	fmt.Printf("  Period high  : $%.2f\n", MaxPrice(prices))
	fmt.Printf("  Period low   : $%.2f\n", MinPrice(prices))
	fmt.Printf("  SMA-20       : $%.2f\n", SimpleMovingAverage(prices, 20))
	fmt.Printf("  SMA-50       : $%.2f\n", SimpleMovingAverage(prices, 50))
	fmt.Printf("  Volatility   : $%.2f  (20-day std dev)\n", Volatility(prices, 20))

	// Train
	PrintHeader("Model Training")
	model := Train(prices)
	fmt.Printf("  Model        : Linear Regression (OLS)\n")
	fmt.Printf("  Slope        : %.4f  ($/day)\n", model.Slope)
	fmt.Printf("  Intercept    : $%.2f\n", model.Intercept)
	fmt.Printf("  R²           : %.4f\n", model.RSquared)

	if model.RSquared < 0.3 {
		fmt.Println("  ⚠  Low R² – the trend is weak; predictions may be unreliable.")
	}

	// Predict 
	PrintHeader(fmt.Sprintf("7-Day Price Forecast for %s", ticker))
	fmt.Printf("  %-12s  %-12s  %-10s\n", "Date", "Predicted", "Change")
	PrintDivider()

	baseIndex := float64(len(prices) - 1)
	lastClose := latest.Close
	currentDate := latest.Date

	for d := 1; d <= forecastDays; d++ {
		// Advance to next trading day (skip weekends)
		currentDate = currentDate.AddDate(0, 0, 1)
		for currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
			currentDate = currentDate.AddDate(0, 0, 1)
		}
		pred := model.Predict(baseIndex + float64(d))
		change := pred - lastClose
		pct := (change / lastClose) * 100
		arrow := "▲"
		if change < 0 {
			arrow = "▼"
		}
		fmt.Printf("  %-12s  $%-11.2f  %s %+.2f (%+.2f%%)\n",
			currentDate.Format("02-01-2006"), pred, arrow, change, pct)
	}

	// Confidence note 
	PrintDivider()
	vol := Volatility(prices, 20)
	fmt.Printf("\n  Confidence band (±1σ): $%.2f – $%.2f\n",
		model.Predict(baseIndex+float64(forecastDays))-vol,
		model.Predict(baseIndex+float64(forecastDays))+vol)
	fmt.Println()
	fmt.Println("  ⚠  Disclaimer: This is a simple linear regression model.")
	fmt.Println("     It does NOT account for market sentiment, news, earnings,")
	fmt.Println("     or any non-linear dynamics. Use for educational purposes only.")
	fmt.Println()
}
