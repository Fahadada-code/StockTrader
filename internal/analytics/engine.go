package analytics

import (
	"math"
	"sync"
)

type RollingMetrics struct {
	Symbol       string
	VWAP         float64
	Volatility   float64
	PriceChange  float64
	VolumeChange float64
}

type Engine struct {
	buffers map[string]*ringBuffer
	mu      sync.Mutex
	window  int
}

type ringBuffer struct {
	prices  []float64
	volumes []float64
	pos     int
	size    int
	full    bool
}

func NewEngine(window int) *Engine {
	return &Engine{
		buffers: make(map[string]*ringBuffer),
		window:  window,
	}
}

func (e *Engine) Process(symbol string, price float64, volume float64) RollingMetrics {
	e.mu.Lock()
	rb, ok := e.buffers[symbol]
	if !ok {
		rb = &ringBuffer{
			prices:  make([]float64, e.window),
			volumes: make([]float64, e.window),
			size:    e.window,
		}
		e.buffers[symbol] = rb
	}
	e.mu.Unlock()

	rb.add(price, volume)

	metrics := RollingMetrics{
		Symbol: symbol,
		VWAP:   rb.computeVWAP(),
	}

	if rb.full || rb.pos > 1 {
		metrics.Volatility = rb.computeVolatility()
		metrics.PriceChange = rb.computeChange()
	}

	return metrics
}

func (rb *ringBuffer) add(p, v float64) {
	rb.prices[rb.pos] = p
	rb.volumes[rb.pos] = v
	rb.pos = (rb.pos + 1) % rb.size
	if rb.pos == 0 {
		rb.full = true
	}
}

func (rb *ringBuffer) computeVWAP() float64 {
	var sumPV, sumV float64
	count := rb.size
	if !rb.full {
		count = rb.pos
	}
	for i := 0; i < count; i++ {
		sumPV += rb.prices[i] * rb.volumes[i]
		sumV += rb.volumes[i]
	}
	if sumV == 0 {
		return 0
	}
	return sumPV / sumV
}

func (rb *ringBuffer) computeVolatility() float64 {
	var sum, sumSq float64
	count := rb.size
	if !rb.full {
		count = rb.pos
	}
	for i := 0; i < count; i++ {
		sum += rb.prices[i]
		sumSq += rb.prices[i] * rb.prices[i]
	}
	mean := sum / float64(count)
	variance := (sumSq / float64(count)) - (mean * mean)
	if variance < 0 {
		return 0
	}
	return math.Sqrt(variance)
}

func (rb *ringBuffer) computeChange() float64 {
	count := rb.size
	if !rb.full {
		count = rb.pos
	}
	if count < 2 {
		return 0
	}
	last := rb.prices[(rb.pos-1+rb.size)%rb.size]
	first := rb.prices[(rb.pos-count+rb.size)%rb.size]
	return ((last - first) / first) * 100
}
