package ingestion

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/Fahadada-code/StockTrader/internal/alphavantage"
	"github.com/Fahadada-code/StockTrader/internal/resilience"
)

type Engine struct {
	client     *alphavantage.Client
	symbols    map[string]time.Time // Last polled time
	interval   time.Duration
	activeSubs func() []string // Callback to get active symbols from Manager/Cache
	cb         *resilience.CircuitBreaker
}

func NewEngine(client *alphavantage.Client, interval time.Duration, activeSubs func() []string, cb *resilience.CircuitBreaker) *Engine {
	return &Engine{
		client:     client,
		symbols:    make(map[string]time.Time),
		interval:   interval,
		activeSubs: activeSubs,
		cb:         cb,
	}
}

func (e *Engine) Run(ctx context.Context, onUpdate func(*alphavantage.QuoteData)) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	backoff := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			symbols := e.activeSubs()
			if len(symbols) == 0 {
				continue
			}

			for _, symbol := range symbols {
				// Simple rate limiting: don't poll same symbol too often
				if last, ok := e.symbols[symbol]; ok && time.Since(last) < e.interval {
					continue
				}

				// Exponential backoff if we hit rate limits elsewhere
				if backoff > 0 {
					wait := time.Duration(1<<uint(backoff)) * time.Second
					wait += time.Duration(rand.Intn(1000)) * time.Millisecond // Jitter
					log.Printf("[Ingestion] Rate limit backoff: waiting %v", wait)
					time.Sleep(wait)
					backoff--
				}

				log.Printf("[Ingestion] Polling %s...", symbol)

				err := e.cb.Execute(func() error {
					quote, err := e.client.GetQuote(symbol)
					if err != nil {
						return err
					}
					e.symbols[symbol] = time.Now()
					onUpdate(quote)
					return nil
				})

				if err != nil {
					log.Printf("[Ingestion] Error polling %s: %v", symbol, err)
					if err.Error() == "rate limit reached or symbol not found" {
						backoff = 5 // Start backoff
					}
					continue
				}

				// Stagger requests to avoid bursting
				time.Sleep(2 * time.Second)
			}
		}
	}
}
