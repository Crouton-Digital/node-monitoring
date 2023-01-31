package server

import (
	"github.com/sirupsen/logrus"
	"node-balancer/internal/metrics"
	"node-balancer/internal/nodemonitoring"
	"node-balancer/internal/server/config"
)

func RunServer() {

	logrus.Info("START APP ")
	config := config.GetServerConfig()

	go nodemonitoring.Run(config)
	metrics.Run(config)

}
