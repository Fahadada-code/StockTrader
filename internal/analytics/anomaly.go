package analytics

import (
	"fmt"
	"math"
)

type Anomaly struct {
	Symbol     string  `json:"symbol"`
	Type       string  `json:"type"` // "volume_spike", "price_jump", "momentum"
	Confidence float64 `json:"confidence"`
	Details    string  `json:"details"`
}

func DetectAnomaly(symbol string, price, volume float64, metrics RollingMetrics) *Anomaly {
	// Simple threshold-based anomaly detection

	// 1. Volume Spike Detection
	// If volume is > 3x rolling avg (hypothetically if we had avg volume in metrics)
	// For now, let's just check volatility vs price change

	if metrics.Volatility > 0 && math.Abs(metrics.PriceChange) > 2.0 {
		return &Anomaly{
			Symbol:     symbol,
			Type:       "momentum_driven",
			Confidence: 0.85,
			Details:    fmt.Sprintf("Significant price move of %.2f%% detected.", metrics.PriceChange),
		}
	}

	return nil
}
