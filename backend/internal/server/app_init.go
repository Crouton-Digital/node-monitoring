package server

import (
	"crypto-exporter/internal/metrics"
	"crypto-exporter/internal/nodemonitoring"
	"crypto-exporter/internal/server/config"
	"github.com/sirupsen/logrus"
)

func RunServer() {

	logrus.Info("START APP ")
	config := config.GetServerConfig()

	go nodemonitoring.Run(config)
	metrics.Run(config)

}
