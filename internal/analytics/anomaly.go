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
	// 1. Price Jump Detection
	if math.Abs(metrics.PriceChange) > 2.0 {
		return &Anomaly{
			Symbol:     symbol,
			Type:       "price_jump",
			Confidence: math.Min(0.95, 0.5+math.Abs(metrics.PriceChange)/10.0),
			Details:    fmt.Sprintf("Sudden momentum shift: %.2f%% price move detected.", metrics.PriceChange),
		}
	}

	// 2. Volume Spike Detection (assuming high volatility correlates with volume spikes in our ring buffer)
	if metrics.Volatility > 1.0 && metrics.VWAP > 0 {
		return &Anomaly{
			Symbol:     symbol,
			Type:       "high_volatility_spike",
			Confidence: 0.8,
			Details:    "Aggressive trading activity detected with elevated volatility.",
		}
	}

	return nil
}
