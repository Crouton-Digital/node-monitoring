# node-balancer #

This project for monitoing EVM nodes and provide configs for traefik http router. 

1. route ws, http trafik between some nodes 
2. monitoring EVM nodes and turn on or turn off from routing 
3. collect and share prometheus metrics 


### RUN PROJECT 

```sh
cd backend 
go mod init node-balancer
go mod tidy

go run app/main.go
```

### Contribution guidelines ###

Configuration app node-balancer

```yaml
server:
    http_port: 8080
    metrics_port: 8082

nodes:
    polygon:
        - label: Polygon ECNG
          url: http://144.76.18.142:36360
          public: false
          proxy_enable: true
          ws_support: true
        - label: Polygon ECNG2
          url: http://65.108.201.189:36360
          public: false
          proxy_enable: true
          ws_support: true
        - label: Polygon RPC
          url: https://polygon-rpc.com
          public: true
          proxy_enable: false
          ws_support: false
        - label: Polygon Infura
          url: https://polygon-mainnet.infura.io/v3/e444d8655de54657b719e041d951aac7
          public: true
          proxy_enable: false
          ws_support: false
    cronos:
        - label: Cronos ECNG
          url: http://142.132.130.196:26659
          public: false
          proxy_enable: true
          ws_support: true
        - label: Cronos
          url: https://evm.cronos.org
          public: true
          proxy_enable: false
          ws_support: false
```

app will be generate configuration for traefik and share use http protocol (example): 

```yaml
http:
  routers:
    polygon_node:
      entryPoints:
        - web
      service: polygon_node
      rule: Host(`polygon-nodes.dex-arbitrage.svc.cluster.local`)
    polygon_node_ws:
      entryPoints:
        - ws
      service: polygon_node_ws
      rule: Host(`polygon-nodes.dex-arbitrage.svc.cluster.local`)
  services:
    polygon_node:
      loadBalancer:
        healthCheck:
          path: /
          port: 36360
          headers:
            Content-Type: "application/json"
          interval: 10s
          timeout: 3s
        servers:
          - url: "http://144.76.18.142:36360"
          - url: "http://65.108.201.189:36360"
        passHostHeader: true
    polygon_node_ws:
      loadBalancer:
        healthCheck:
          path: /
          port: 36360
          headers:
            Content-Type: "application/json"
          interval: 10s
          timeout: 3s
        servers:
          - url: "http://144.76.18.142:36361"
          - url: "http://65.108.201.189:36361"
        passHostHeader: true
```

### How to analyze EVN nodes?  

### List supported prometheus metrics 