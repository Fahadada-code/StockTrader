package alphavantage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

type Client struct {
	apiKey  string
	baseURL string
	cache   map[string]cacheEntry
	mu      sync.RWMutex
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://www.alphavantage.co/query",
		cache:   make(map[string]cacheEntry),
	}
}

// Internal structure to match Alpha Vantage API response
type avGlobalQuoteResponse struct {
	GlobalQuote avQuoteData `json:"Global Quote"`
}

type avQuoteData struct {
	Symbol           string `json:"01. symbol"`
	Open             string `json:"02. open"`
	High             string `json:"03. high"`
	Low              string `json:"04. low"`
	Price            string `json:"05. price"`
	Volume           string `json:"06. volume"`
	LatestTradingDay string `json:"07. latest trading day"`
	PreviousClose    string `json:"08. previous close"`
	Change           string `json:"09. change"`
	ChangePercent    string `json:"10. change percent"`
}

// Clean structure for our internal API and Frontend
type QuoteData struct {
	Symbol           string
	Open             string
	High             string
	Low              string
	Price            string
	Volume           string
	LatestTradingDay string
	PreviousClose    string
	Change           string
	ChangePercent    string
}

func (c *Client) GetQuote(symbol string) (*QuoteData, error) {
	cacheKey := "quote:" + symbol
	c.mu.RLock()
	entry, found := c.cache[cacheKey]
	c.mu.RUnlock()

	if found && time.Now().Before(entry.expiresAt) {
		return entry.data.(*QuoteData), nil
	}

	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", c.baseURL, symbol, c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result avGlobalQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.GlobalQuote.Symbol == "" {
		// Alpha Vantage returns a 200 with an error message in the JSON when limited
		return nil, fmt.Errorf("rate limit reached or symbol not found")
	}

	quote := &QuoteData{
		Symbol:           result.GlobalQuote.Symbol,
		Open:             result.GlobalQuote.Open,
		High:             result.GlobalQuote.High,
		Low:              result.GlobalQuote.Low,
		Price:            result.GlobalQuote.Price,
		Volume:           result.GlobalQuote.Volume,
		LatestTradingDay: result.GlobalQuote.LatestTradingDay,
		PreviousClose:    result.GlobalQuote.PreviousClose,
		Change:           result.GlobalQuote.Change,
		ChangePercent:    result.GlobalQuote.ChangePercent,
	}

	c.mu.Lock()
	c.cache[cacheKey] = cacheEntry{
		data:      quote,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
	c.mu.Unlock()

	return quote, nil
}

// Internal structure for Daily History
type avTimeSeriesDailyResponse struct {
	TimeSeries map[string]avDailyData `json:"Time Series (Daily)"`
}

type avDailyData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

// Clean structure for Frontend
type DailyData struct {
	Open   string
	High   string
	Low    string
	Close  string
	Volume string
}

func (c *Client) GetDailyHistory(symbol string) (map[string]DailyData, error) {
	cacheKey := "history:" + symbol
	c.mu.RLock()
	entry, found := c.cache[cacheKey]
	c.mu.RUnlock()

	if found && time.Now().Before(entry.expiresAt) {
		return entry.data.(map[string]DailyData), nil
	}

	url := fmt.Sprintf("%s?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", c.baseURL, symbol, c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result avTimeSeriesDailyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.TimeSeries == nil {
		return nil, fmt.Errorf("rate limit reached or history not found")
	}

	cleanHistory := make(map[string]DailyData)
	for date, data := range result.TimeSeries {
		cleanHistory[date] = DailyData{
			Open:   data.Open,
			High:   data.High,
			Low:    data.Low,
			Close:  data.Close,
			Volume: data.Volume,
		}
	}

	c.mu.Lock()
	c.cache[cacheKey] = cacheEntry{
		data:      cleanHistory,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
	c.mu.Unlock()

	return cleanHistory, nil
}
