package trafik_config

import (
	"fmt"
	"log"
	"net/http"
	"node-balancer/internal/nodemonitoring"
	"node-balancer/internal/server/config"

	"golang.org/x/exp/slices"
	yaml3 "gopkg.in/yaml.v3"
)

type TrafikConfig struct {
	HTTP HTTP `yaml:"http"`
}

type Router struct {
	EntryPoints []string `yaml:"entryPoints"`
	Service     string   `yaml:"service"`
	Rule        string   `yaml:"rule"`
}

type HealthCheck struct {
	Path     string            `yaml:"path"`
	Port     int               `yaml:"port"`
	Headers  map[string]string `yaml:"headers"`
	Interval string            `yaml:"interval"`
	Timeout  string            `yaml:"timeout"`
}
type Server struct {
	URL string `yaml:"url"`
}
type LoadBalancer struct {
	// HealthCheck    HealthCheck `yaml:"healthCheck"`
	Servers        []Server `yaml:"servers"`
	PassHostHeader bool     `yaml:"passHostHeader"`
}

type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

type HTTP struct {
	Routers  map[string]Router  `yaml:"routers"`
	Services map[string]Service `yaml:"services"`
}

func GenerateConfig() TrafikConfig {

	tfConfig := TrafikConfig{HTTP: HTTP{
		Routers:  map[string]Router{},
		Services: map[string]Service{},
	}}

	for network, netConfig := range config.Config.NetworksConfig {
		tfConfig.HTTP.Routers[network+"_node"] = Router{
			EntryPoints: []string{"web"},
			Service:     network + "_node",
			Rule:        "Host(`" + netConfig.Domain + "`)",
		}
		tfConfig.HTTP.Routers[network+"_node_ws"] = Router{
			EntryPoints: []string{"ws"},
			Service:     network + "_node_ws",
			Rule:        "Host(`" + netConfig.Domain + "`)",
		}
	}

	for network, netConfig := range config.Config.NetworksConfig {
		httpServers := []Server{}
		wsServers := []Server{}

		enabledNodes := nodemonitoring.EnabledNodes(network)

		for index, node := range netConfig.Nodes {
			if node.AllowRouting && slices.Contains(enabledNodes, index) {
				if node.Url != "" {
					httpServers = append(httpServers, Server{URL: node.Url})
				}

				if node.WsUrl != "" {
					wsServers = append(wsServers, Server{URL: node.WsUrl})
				}
			}
		}

		tfConfig.HTTP.Services[network+"_node"] = Service{
			LoadBalancer: LoadBalancer{
				// HealthCheck: HealthCheck{
				// 	Path: "/",
				// 	Port: 36360,
				// 	Headers: map[string]string{
				// 		"Content-Type": "application/json",
				// 	},
				// 	Interval: "10s",
				// 	Timeout:  "3s",
				// },
				Servers:        httpServers,
				PassHostHeader: true,
			},
		}
		tfConfig.HTTP.Services[network+"_node_ws"] = Service{
			LoadBalancer: LoadBalancer{
				Servers:        wsServers,
				PassHostHeader: true,
			},
		}
	}

	return tfConfig
}

func HandleConfig(w http.ResponseWriter, r *http.Request) {
	tfConfig := GenerateConfig()

	configBytes, err := yaml3.Marshal(tfConfig)
	if err != nil {
		log.Fatalf("can not marshal trafik config: %v", err)
	}
	w.Header().Del("Transfer-Encoding")
	w.Header().Set("Content-Type", "text/yaml")
	w.Header().Set("Content-Length", fmt.Sprint(len(configBytes)))
	w.WriteHeader(http.StatusOK)
	w.Write(configBytes)
}
