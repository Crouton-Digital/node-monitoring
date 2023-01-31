package server

import (
	"fmt"
	"net/http"
	"node-balancer/internal/nodemonitoring"
	"node-balancer/internal/server/config"
	"node-balancer/internal/trafik_config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func RunServer() {
	config := config.GetServerConfig()
	logrus.Info("START APP ")
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetLevel(logrus.DebugLevel)
	go nodemonitoring.Run(config)
	StartRouter(config)
}

func StartRouter(config config.Config) {
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/health", http.HandlerFunc(handleHealthRequest))
	http.Handle("/traefik-route.cfg", http.HandlerFunc(trafik_config.HandleConfig))

	logrus.Infof("Serving HTTP on port %v", config.ServerConfig.HttpPort)

	http.ListenAndServe(":"+config.ServerConfig.HttpPort, nil)
}

func handleHealthRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}
