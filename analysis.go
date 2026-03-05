package main

import "math"



// SimpleMovingAverage returns the SMA of the last `window` closing prices.
func SimpleMovingAverage(prices []float64, window int) float64 {
	if len(prices) < window {
		window = len(prices)
	}
	sum := 0.0
	for _, p := range prices[len(prices)-window:] {
		sum += p
	}
	return sum / float64(window)
}

// Volatility returns the standard deviation of the last `window` closing prices.
func Volatility(prices []float64, window int) float64 {
	if len(prices) < window {
		window = len(prices)
	}
	subset := prices[len(prices)-window:]
	mean := 0.0
	for _, p := range subset {
		mean += p
	}
	mean /= float64(len(subset))

	variance := 0.0
	for _, p := range subset {
		diff := p - mean
		variance += diff * diff
	}
	variance /= float64(len(subset))
	return math.Sqrt(variance)
}

// MaxPrice returns the highest value in the slice.
func MaxPrice(prices []float64) float64 {
	m := prices[0]
	for _, p := range prices[1:] {
		if p > m {
			m = p
		}
	}
	return m
}

// MinPrice returns the lowest value in the slice.
func MinPrice(prices []float64) float64 {
	m := prices[0]
	for _, p := range prices[1:] {
		if p < m {
			m = p
		}
	}
	return m
}
