package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Fahadada-code/StockTrader/internal/alphavantage"
	"github.com/Fahadada-code/StockTrader/internal/analytics"
	"github.com/Fahadada-code/StockTrader/internal/cache"
	"github.com/Fahadada-code/StockTrader/internal/db"
	"github.com/Fahadada-code/StockTrader/internal/ingestion"
	"github.com/Fahadada-code/StockTrader/internal/metrics"
	"github.com/Fahadada-code/StockTrader/internal/websocket"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()

	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		log.Fatal("ALPHA_VANTAGE_API_KEY is not set")
	}

	// 1. Initialize DB & Cache
	pgURL := os.Getenv("DB_URL")
	if pgURL == "" {
		pgURL = "postgres://postgres:postgres@127.0.0.1:5433/stocktrader?sslmode=disable"
	}
	pg, err := db.NewPostgresDB(pgURL)
	if err != nil {
		log.Printf("Warning: Database connection failed (%v). Entering in-memory mode for persistence.", err)
	}

	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6380"
	}

	// Check Redis connectivity before initializing the full client to avoid noisy background dials
	var redisAvailable bool
	rc := cache.NewRedisCache(redisAddr)
	ctxCheck, cancelCheck := context.WithTimeout(context.Background(), 1*time.Second)
	if err := rc.Client.Ping(ctxCheck).Err(); err == nil {
		redisAvailable = true
	} else {
		log.Printf("Warning: Redis connection failed (%v). Entering in-memory mode for caching.", err)
		rc.Client.Close() // Close to stop background connection attempts
	}
	cancelCheck()

	// 2. Initialize Components
	avClient := alphavantage.NewClient(apiKey)
	wsManager := websocket.NewManager()
	analyticsEngine := analytics.NewEngine(50) // 50-point rolling window

	activeSymbolsFunc := func() []string {
		if !redisAvailable {
			return []string{"AAPL", "TSLA", "MSFT"} // default set
		}
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		symbols, err := rc.GetHotSymbols(ctx, 10)
		if err != nil {
			log.Printf("Warning: Redis query failed. Falling back to defaults. Error: %v", err)
			redisAvailable = false
			return []string{"AAPL", "TSLA", "MSFT"}
		}
		if len(symbols) == 0 {
			return []string{"AAPL", "TSLA", "MSFT"}
		}
		// Clean the keys (subs:AAPL -> AAPL)
		for i, s := range symbols {
			if len(s) > 5 {
				symbols[i] = s[5:]
			}
		}
		return symbols
	}

	ingestionEngine := ingestion.NewEngine(avClient, 30*time.Second, activeSymbolsFunc)

	// 3. Start Background Routines
	go wsManager.Run()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ingestionEngine.Run(ctx, func(quote *alphavantage.QuoteData) {
		price, _ := strconv.ParseFloat(quote.Price, 64)
		volume, _ := strconv.ParseInt(quote.Volume, 10, 64)

		// A. Analytics Processing
		m := analyticsEngine.Process(quote.Symbol, price, float64(volume))
		metrics.UpdatesProcessed.WithLabelValues(quote.Symbol).Inc()

		// B. Anomaly Detection
		if anomaly := analytics.DetectAnomaly(quote.Symbol, price, float64(volume), m); anomaly != nil {
			metrics.AnomaliesDetected.WithLabelValues(quote.Symbol, anomaly.Type).Inc()
			wsManager.Broadcast(websocket.Message{
				Symbol: quote.Symbol,
				Type:   "anomaly",
				Data:   anomaly,
			})
		}

		// C. Snapshot Persistence
		if pg != nil {
			start := time.Now()
			pg.SaveQuote(quote.Symbol, price, volume)
			metrics.DatabaseLatency.Observe(time.Since(start).Seconds())
		}

		// D. Real-Time Distribution
		wsManager.Broadcast(websocket.Message{
			Symbol: quote.Symbol,
			Type:   "price",
			Data:   quote,
		})
	})

	// 4. HTTP Handlers
	enableCORS := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
			next(w, r)
		}
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(wsManager, w, r)
		metrics.ActiveConnections.Inc()
	})

	http.HandleFunc("/api/quote", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			http.Error(w, "symbol is required", http.StatusBadRequest)
			return
		}
		quote, err := avClient.GetQuote(symbol)
		if err != nil {
			if err.Error() == "rate limit reached or symbol not found" {
				http.Error(w, err.Error(), http.StatusTooManyRequests)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(quote)
	}))

	http.HandleFunc("/api/history", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			http.Error(w, "symbol is required", http.StatusBadRequest)
			return
		}
		history, err := avClient.GetDailyHistory(symbol)
		if err != nil {
			if err.Error() == "rate limit reached or history not found" {
				http.Error(w, err.Error(), http.StatusTooManyRequests)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(history)
	}))

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/api/health", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"up"}`))
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("StockTrader Pro Backend starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
