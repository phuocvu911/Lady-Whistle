package main

// LinearModel holds slope (m) and intercept (b) for y = m*x + b.
type LinearModel struct {
	Slope     float64
	Intercept float64
	RSquared  float64
}

// Train fits a least-squares linear regression on the data.
// x = day index (0, 1, 2, …), y = closing price.
func Train(prices []float64) LinearModel {
	n := float64(len(prices))
	var sumX, sumY, sumXY, sumX2 float64

	for i, y := range prices {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// slope  m = (n*Σxy - Σx*Σy) / (n*Σx² - (Σx)²)
	denom := n*sumX2 - sumX*sumX
	m := (n*sumXY - sumX*sumY) / denom
	// intercept b = (Σy - m*Σx) / n
	b := (sumY - m*sumX) / n

	// R² (coefficient of determination)
	meanY := sumY / n
	var ssTot, ssRes float64
	for i, y := range prices {
		pred := m*float64(i) + b
		ssRes += (y - pred) * (y - pred)
		ssTot += (y - meanY) * (y - meanY)
	}
	r2 := 1.0
	if ssTot > 0 {
		r2 = 1.0 - ssRes/ssTot
	}

	return LinearModel{Slope: m, Intercept: b, RSquared: r2}
}

// Predict returns the predicted price for a given day index.
func (lm LinearModel) Predict(dayIndex float64) float64 {
	return lm.Slope*dayIndex + lm.Intercept
}
