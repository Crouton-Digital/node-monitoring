package config

import (
	"os"

	"github.com/sirupsen/logrus"
	yaml3 "gopkg.in/yaml.v3"
)

type AppConfig struct {
	ServerConfig   Server             `yaml:"server"`
	NetworksConfig map[string]Network `yaml:"networks"`
}

type Server struct {
	HttpPort   string `yaml:"http_port"`
	DebugLevel string `yaml:"debug"`
}

type Network struct {
	Domain string       `yaml:"domain"`
	Rules  NetworkRules `yaml:"rules"`
	Nodes  []Node       `yaml:"nodes"`
}

type NetworkRules struct {
	MaxBlockDelay   int `yaml:"max_block_delay"`
	MaxTimeDelaySec int `yaml:"max_time_delay_sec"`
	RoutingNodesMin int `yaml:"routing_nodes_min"`
	RoutingNodesMax int `yaml:"routing_nodes_max"`
}

type Node struct {
	Name         string `yaml:"label"`
	Url          string `yaml:"url"`
	WsUrl        string `yaml:"ws_url"`
	Public       bool   `yaml:"public"`
	AllowRouting bool   `yaml:"allow_routing"`
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

	//logrus.Info(Config)
	for key, network_nodes := range Config.NetworksConfig {
		logrus.Infof("======== %v ========", key)
		for _, network_node := range network_nodes.Nodes {
			logrus.Infof("%v | %v %v Public: %v", key, network_node.Name, network_node.Url, network_node.Public)
		}
	}
}
