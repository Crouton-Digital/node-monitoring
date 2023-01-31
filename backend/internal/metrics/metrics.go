package metrics

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"node-balancer/internal/server/config"
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

func Run(config config.Config) {

	handlerHealth := http.HandlerFunc(handleHealthRequest)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/health", handlerHealth)
	http.Handle("/traefik-route.cfg", handlerHealth)

	http.ListenAndServe(":"+config.ServerConfig.HttpPort, nil)
}

func handleHealthRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		logrus.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
