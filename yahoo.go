package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

// YahooResponse maps the JSON returned by Yahoo Finance's chart API.
type YahooResponse struct {
	Chart struct {
		Result []struct {
			Timestamp  []int64 `json:"timestamp"` //struct tag to tell the unmarshaler how to map JSON data to the struct field.
			Indicators struct {
				Quote []struct {
					Close  []*float64 `json:"close"`
					Open   []*float64 `json:"open"`
					High   []*float64 `json:"high"`
					Low    []*float64 `json:"low"`
					Volume []*int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

// StockData holds one day of stock data.
type StockData struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

// ─────────────────────────────────────────────
// Fetch historical data from Yahoo Finance
// ─────────────────────────────────────────────

// FetchStockData downloads daily OHLCV data for the given ticker
// over the specified number of calendar days into the past.
func FetchStockData(ticker string, days int) ([]StockData, error) {
	now := time.Now()
	from := now.AddDate(0, 0, -days)

	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?period1=%d&period2=%d&interval=1d",
		ticker, from.Unix(), now.Unix(),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	var yahoo YahooResponse
	if err := json.Unmarshal(body, &yahoo); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	if yahoo.Chart.Error != nil {
		return nil, fmt.Errorf("yahoo API error: %s – %s",
			yahoo.Chart.Error.Code, yahoo.Chart.Error.Description)
	}

	if len(yahoo.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data returned for ticker %q", ticker)
	}

	result := yahoo.Chart.Result[0]
	quotes := result.Indicators.Quote[0]

	var data []StockData
	for i, ts := range result.Timestamp {
		// Skip entries with nil close price (market holidays, etc.)
		if quotes.Close[i] == nil {
			continue
		}
		sd := StockData{
			Date:  time.Unix(ts, 0),
			Close: *quotes.Close[i],
		}
		if quotes.Open[i] != nil {
			sd.Open = *quotes.Open[i]
		}
		if quotes.High[i] != nil {
			sd.High = *quotes.High[i]
		}
		if quotes.Low[i] != nil {
			sd.Low = *quotes.Low[i]
		}
		if quotes.Volume[i] != nil {
			sd.Volume = *quotes.Volume[i]
		}
		data = append(data, sd)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Date.Before(data[j].Date)
	})

	return data, nil
}
