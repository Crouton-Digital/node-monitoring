package config

import (
	"github.com/sirupsen/logrus"
	yaml3 "gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServerConfig  Server             `yaml:"server"`
	DomainsConfig map[string]Domains `yaml:"domains"`
	RpcConfig     map[string][]Node  `yaml:"nodes"`
}

type Server struct {
	HttpPort    string `yaml:"http_port"`
	MetricsPort string `yaml:"metrics_port"`
	DebugLevel  string `yaml:"debug"`
}

type Domains struct {
	Url string `yaml:"url"`
}

type Node struct {
	Name     string `yaml:"label"`
	Url      string `yaml:"url"`
	Public   bool   `yaml:"public"`
	RouteUrl string `yaml:"route_url"`
}

func GetServerConfig() Config {

	env := os.Getenv("ENV")

	if env == "" {
		env = "local"
	}

	data, err := os.ReadFile("config/" + env + ".yml")
	if err != nil {
		logrus.Errorf("Failed to read config: %v", err)
		os.Exit(1)
	}
	var config Config
	err = yaml3.Unmarshal(data, &config)
	if err != nil {
		logrus.Errorf("Failed to parse config: %v", err)
		os.Exit(1)
	}

	logrus.Info(config.DomainsConfig)
	for key, network_nodes := range config.RpcConfig {
		logrus.Infof("======== %v ========", key)
		for _, network_node := range network_nodes {
			logrus.Infof("%v", config.DomainsConfig[key].Url)
			logrus.Infof("%v | %v %v Public: %v", key, network_node.Name, network_node.Url, network_node.Public)
		}
	}

	return config
}
