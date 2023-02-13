package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	timeBuckets = []float64{.010, .025, .05, .1, .25, .5, 1, 2, 3, 4, 5, 7, 10}

	responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "node_balancer_response",
			Help:    "Response time, partitioned by HTTP status code",
			Buckets: timeBuckets,
		},
		[]string{"network", "node_name", "error"},
	)

	blockNum = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_balancer_block",
		Help: "Last block height",
	}, []string{"network", "node_name"})

	blockTimeAgo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_balancer_block_time_ago",
		Help: "Last block happened X seconds ago",
	}, []string{"network", "node_name"})

	blockDelayFromBest = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_balancer_block_delay_from_best",
		Help: "Diff in blocks between last block and best last block",
	}, []string{"network", "node_name"})

	// node_balancer_avability
	// +node_balancer_http_code_response
	// node_balancer_block
	// node_balancer_block_time_ago
	// node_balancer_block_delay_from_best
	// node_balancer_allow_routing
)

func init() {
	prometheus.MustRegister(responseTime)
	prometheus.MustRegister(blockNum)
	prometheus.MustRegister(blockTimeAgo)
	prometheus.MustRegister(blockDelayFromBest)
}

func ResponseTime(network, nodeName, err string, start time.Time) {
	dur := time.Since(start)
	responseTime.WithLabelValues(network, nodeName, err).Observe(dur.Seconds())
}

func BlockNum(network, nodeName string, blockNumber int64) {
	blockNum.WithLabelValues(network, nodeName).Set(float64(blockNumber))
}

func BlockTimeAgo(network, nodeName string, blockTime time.Time) {
	blockTimeAgo.WithLabelValues(network, nodeName).Set(time.Since(blockTime).Seconds())
}

func BlockDelayFromBest(network, nodeName string, diff int64) {
	blockDelayFromBest.WithLabelValues(network, nodeName).Set(float64(diff))
}
