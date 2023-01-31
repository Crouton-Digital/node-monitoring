package server

import (
	"github.com/sirupsen/logrus"
	"node-balancer/internal/metrics"
	"node-balancer/internal/nodemonitoring"
	"node-balancer/internal/server/config"
)

func RunServer() {

	config := config.GetServerConfig()
	logrus.Info("START APP ")
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetLevel(logrus.DebugLevel)
	go nodemonitoring.Run(config)
	metrics.Run(config)

}
