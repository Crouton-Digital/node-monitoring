package nodemonitoring

import (
	"context"
	"node-balancer/internal/node_rating"
	"node-balancer/internal/server/config"
	"node-balancer/internal/utils"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var (
	// stores indexes, e.g. for polygon enabled nodes are #0 and #3
	enabledNodes   = map[string][]int{}
	mu             = sync.RWMutex{}
	WorkerChannels = make(map[string](chan struct{}))
)

func Run() {
	for network := range config.Config.NetworksConfig {
		startScheduler(network)
	}

	for range time.Tick(time.Millisecond * 2000) {
		logrus.Infof("======= %v =======", time.Now().Format("2006-01-02 15:04:05"))
		for network := range config.Config.NetworksConfig {
			WorkerChannels[network] <- struct{}{}
		}
	}
}

func startScheduler(network string) {
	WorkerChannels[network] = make(chan struct{})
	go worker(network, WorkerChannels[network])
}

func worker(network string, ch <-chan struct{}) {
	for range ch {
		monitorNetwork(network)
	}
}

func IsEnabled(network string, index int) bool {
	enabled := EnabledNodes(network)
	return slices.Contains(enabled, index)
}

func EnabledNodes(network string) []int {
	mu.RLock()
	defer mu.RUnlock()
	return enabledNodes[network]
}

func setEnabledNodes(network string, nodeIndexes []int) {
	mu.Lock()
	defer mu.Unlock()
	enabledNodes[network] = nodeIndexes
}

type monitoredNode struct {
	Index         int
	LastBlock     int64
	LastBlockTime time.Time
	BlockDelay    int64 // number of blocks since best last block
	Error         error
}

func monitorNetwork(network string) {
	netConfig := config.Config.NetworksConfig[network]
	var bestlastBlock int64

	monitoredNodes := utils.ParallelMap(netConfig.Nodes, func(i int, node config.Node) monitoredNode {
		//logrus.Infof("    %s.%d %s - Checking last block", network, i, node.Name)
		lastBlock, lastBlockTime, err := getLastKnowBlock(node)
		logrus.Infof("    %s.%d %s - block %d | %v ago | %v", network, i, node.Name, lastBlock, time.Since(lastBlockTime), err)

		return monitoredNode{Index: i, LastBlock: lastBlock, LastBlockTime: lastBlockTime, Error: err}
	})

	for _, mnode := range monitoredNodes {
		if mnode.Error == nil && mnode.LastBlock > bestlastBlock {
			bestlastBlock = mnode.LastBlock
		}
	}

	// Calculate BlockDelay
	for i := range netConfig.Nodes {
		monitoredNodes[i].BlockDelay = bestlastBlock - monitoredNodes[i].LastBlock
		node_rating.AddRating(network, i, monitoredNodes[i].BlockDelay, monitoredNodes[i].Error)
	}

	//Print info for all nodes:
	// sortedNodes := node_rating.NodesSortedByRating(network)
	// for idx, node := range sortedNodes {
	// 	logrus.Infof("%s %d |%s", network, idx, node.Label())
	// }

	topNodes := node_rating.RoutableNodesWithBestRatings(network)

	topNodesStr := []string{}
	for _, node := range topNodes {
		topNodesStr = append(topNodesStr, node.Label())
	}
	logrus.Infof("%s | best block: %d | top routable nodes: %+v", network, bestlastBlock, strings.Join(topNodesStr, ", "))

	netEnabledNodes := []int{}
	for _, node := range topNodes {
		netEnabledNodes = append(netEnabledNodes, node.Index)
	}

	setEnabledNodes(network, netEnabledNodes)
}

// func printBlockNumber(s config.Node) {
// 	start := time.Now()
// 	_, blockNum := getBlockNumber(s)
// 	spent := time.Since(start)

// 	getLastBlock, getLastBlockTime, err := getLastKnowBlock(s)
// 	if err != nil {

// 	} else {
// 		metrics.OpsBlockHight.WithLabelValues("polygon", s.Name, s.Url, strconv.FormatBool(s.Public)).Set(float64(getLastBlock))
// 	}

// 	w := new(tabwriter.Writer)
// 	w.Init(os.Stdout, 15, 20, 0, '\t', 0)
// 	defer w.Flush()
// 	fmt.Fprintf(w, "\n %v\t %v\t %v\t %v\t %v\t %s ", start.Format("15:04:05.99999"), spent, blockNum, getLastBlock, getLastBlockTime.Format("15:04:05.99999"), s.Name)
// }

// func getBlockNumber(s config.Node) (error, uint64) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	client, err := ethclient.DialContext(ctx, s.Url)
// 	if err != nil {
// 		logrus.Errorf("\n Error connect to: %v Error: %v", s.Url, err)
// 	}

// 	header, err := client.BlockNumber(ctx)
// 	if err != nil {
// 		logrus.Errorf("\nGet block num. Error: %v", err)
// 		return err, 0
// 	} else {
// 		return nil, header
// 	}
// }

func getLastKnowBlock(s config.Node) (int64, time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, s.Url)
	if err != nil {
		logrus.Errorf("\nError connect to: %v Error: %v", s.Url, err)
		return 0, time.Unix(0, 0), err
	}

	latesHeader, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		logrus.Errorf("Get last know block num. %v", err)
		return 0, time.Unix(0, 0), err
	}
	blockTime := time.Unix(int64(latesHeader.Time), 0)
	return latesHeader.Number.Int64(), blockTime, nil
}
