package trafik_config

import (
	"log"
	"net/http"

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
	HealthCheck    HealthCheck `yaml:"healthCheck"`
	Servers        []Server    `yaml:"servers"`
	PassHostHeader bool        `yaml:"passHostHeader"`
}

type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

type HTTP struct {
	Routers  map[string]Router  `yaml:"routers"`
	Services map[string]Service `yaml:"services"`
}

func GenerateConfig() TrafikConfig {

	config := TrafikConfig{
		HTTP: HTTP{
			Routers: map[string]Router{
				"polygon_node": {
					EntryPoints: []string{"web"},
					Service:     "polygon_node",
					Rule:        "Host(`polygon-nodes.dex-arbitrage.svc.cluster.local`)",
				},
				"polygon_node_ws": {
					EntryPoints: []string{"ws"},
					Service:     "polygon_node_ws",
					Rule:        "Host(`polygon-nodes.dex-arbitrage.svc.cluster.local`)",
				},
			},
			Services: map[string]Service{
				"polygon_node": {
					LoadBalancer: LoadBalancer{
						HealthCheck: HealthCheck{
							Path: "/",
							Port: 36360,
							Headers: map[string]string{
								"Content-Type": "application/json",
							},
							Interval: "10s",
							Timeout:  "3s",
						},
						Servers: []Server{
							{URL: "http://144.76.18.142:36360"},
							{URL: "http://65.108.201.189:36360"},
						},
					},
				},
				"polygon_node_ws": {
					LoadBalancer: LoadBalancer{
						HealthCheck: HealthCheck{
							Path: "/",
							Port: 36360,
							Headers: map[string]string{
								"Content-Type": "application/json",
							},
							Interval: "10s",
							Timeout:  "3s",
						},
						Servers: []Server{
							{URL: "http://144.76.18.142:36361"},
							{URL: "http://65.108.201.189:36361"},
						},
						PassHostHeader: true,
					},
				},
			},
		},
	}

	return config
}

func HandleConfig(w http.ResponseWriter, r *http.Request) {
	tfConfig := GenerateConfig()

	configBytes, err := yaml3.Marshal(tfConfig)
	if err != nil {
		log.Fatalf("can not marshal trafik config: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/yaml")
	w.Write(configBytes)
}
