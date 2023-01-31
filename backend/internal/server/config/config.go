package config

import (
	"github.com/sirupsen/logrus"
	yaml3 "gopkg.in/yaml.v3"
	"os"
	"reflect"
)

type Config struct {
	RpcConfig map[string][]Node `yaml:"nodes"`
}

//type Server

type Node struct {
	Name   string `yaml:"label"`
	Type   string
	Url    string `yaml:"url"`
	Public bool   `yaml:"public"`
}

func GetServerConfig() Config {

	data, err := os.ReadFile(".env")
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

	config_group := reflect.ValueOf(config.RpcConfig).MapKeys()
	//logrus.Println(config.RpcConfig["polygon"][0].Name)
	for _, s := range config_group {
		for i, v := range config.RpcConfig[s.String()] {
			config.RpcConfig[s.String()][i].Type = s.String()
			logrus.Infof("%v | %v %v Public: %v", v.Type, v.Name, v.Url, v.Public)
		}
	}

	return config
}
