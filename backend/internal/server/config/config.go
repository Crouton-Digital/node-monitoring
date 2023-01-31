package config

import (
	"os"

	"github.com/sirupsen/logrus"
	yaml3 "gopkg.in/yaml.v3"
)

type AppConfig struct {
	ServerConfig  Server             `yaml:"server"`
	DomainsConfig map[string]Domains `yaml:"domains"`
	Nodes         map[string][]Node  `yaml:"nodes"`
}

type Server struct {
	HttpPort   string `yaml:"http_port"`
	DebugLevel string `yaml:"debug"`
}

type Domains struct {
	Url string `yaml:"url"`
}

type Node struct {
	Name         string `yaml:"label"`
	Url          string `yaml:"url"`
	Public       bool   `yaml:"public"`
	ProxyEnabled bool   `yaml:"proxy_enabled"`
	WsSupport    bool   `yaml:"ws_support"`
}

var (
	Config AppConfig
)

func LoadServerConfig() {

	env := os.Getenv("ENV")

	if env == "" {
		env = "local"
	}

	data, err := os.ReadFile("config/" + env + ".yml")
	if err != nil {
		logrus.Errorf("Failed to read config: %v", err)
		os.Exit(1)
	}
	err = yaml3.Unmarshal(data, &Config)
	if err != nil {
		logrus.Errorf("Failed to parse config: %v", err)
		os.Exit(1)
	}

	logrus.Info(Config.DomainsConfig)
	for key, network_nodes := range Config.Nodes {
		logrus.Infof("======== %v ========", key)
		for _, network_node := range network_nodes {
			logrus.Infof("%v", Config.DomainsConfig[key].Url)
			logrus.Infof("%v | %v %v Public: %v", key, network_node.Name, network_node.Url, network_node.Public)
		}
	}
}
