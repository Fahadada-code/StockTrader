package ingestion

import (
	"context"
	"log"
	"time"

	"github.com/Fahadada-code/StockTrader/internal/alphavantage"
	"github.com/Fahadada-code/StockTrader/internal/db"
)

type ReplayEngine struct {
	db *db.PostgresDB
}

func NewReplayEngine(pg *db.PostgresDB) *ReplayEngine {
	return &ReplayEngine{db: pg}
}

func (r *ReplayEngine) Replay(ctx context.Context, symbol string, speed float64, onUpdate func(*alphavantage.QuoteData)) {
	rows, err := r.db.GetHistoricalData(symbol, 1000)
	if err != nil {
		log.Printf("Replay error: %v", err)
		return
	}
	defer rows.Close()

	var data []alphavantage.QuoteData
	for rows.Next() {
		var q alphavantage.QuoteData
		var price float64
		var volume int64
		var ts time.Time
		if err := rows.Scan(&price, &volume, &ts); err != nil {
			continue
		}
		// Convert back to string for the existing QuoteData structure (for simplicity)
		q.Symbol = symbol
		q.Price = time.Now().Format("2006-01-02 15:04:05") // Dummy TS or price formatting
		// ... realistically we'd use a more generic data structure
		data = append(data, q)
	}

	for _, q := range data {
		select {
		case <-ctx.Done():
			return
		default:
			onUpdate(&q)
			time.Sleep(time.Duration(float64(time.Second) / speed))
		}
	}
}
