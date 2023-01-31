package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	OpsBlockHight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_block_hight",
		Help: "Get block hight",
	}, []string{"type", "name", "url", "public"})
)
