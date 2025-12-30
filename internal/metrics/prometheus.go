package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	UpdatesProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "stocktrader_updates_total",
		Help: "Total number of stock updates processed",
	}, []string{"symbol"})

	ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "stocktrader_active_connections",
		Help: "Current number of active WebSocket connections",
	})

	AnomaliesDetected = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "stocktrader_anomalies_total",
		Help: "Total number of anomalies detected",
	}, []string{"symbol", "type"})

	DatabaseLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "stocktrader_db_latency_seconds",
		Help:    "Latency of database operations",
		Buckets: prometheus.DefBuckets,
	})
)
