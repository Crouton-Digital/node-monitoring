package metrics

import (
	"crypto-exporter/internal/server/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

//func recordMetrics() {
//	go func() {
//		for {
//			opsProcessed.Add(1)
//			opsCpuUsage.Set(1)
//			time.Sleep(2 * time.Second)
//		}
//	}()
//}

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

func Run(config config.Config) {
	//recordMetrics()
	//prometheus.MustRegister(OpsBlockHight)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
