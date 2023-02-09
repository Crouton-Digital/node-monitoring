package node_rating

import (
	"fmt"
	"node-balancer/internal/server/config"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	nodeRatings = map[string](map[int][]int64){}
	mu          = sync.RWMutex{}
)

type NodeWithRating struct {
	Network string
	Index   int
	Rating  int64
}

func (n NodeWithRating) Label() string {
	node := config.Config.NetworksConfig[n.Network].Nodes[n.Index]

	return fmt.Sprintf("%s (r%d)", node.Name, n.Rating)
}

func AddRating(network string, index int, blockDelayBlocks int64, err error) {
	rating := blockDelayBlocks
	if err != nil {
		rating = config.Config.NodeRating.ErrorRating
	}

	logrus.Infof("AddRating(%s.%d) Lock", network, index)
	mu.Lock()
	defer mu.Unlock()
	logrus.Infof("AddRating(%s.%d) Locked", network, index)

	if nodeRatings[network] == nil {
		nodeRatings[network] = map[int][]int64{}
	}
	ratings, ok := nodeRatings[network][index]
	if !ok {
		nodeRatings[network][index] = []int64{rating}
		return
	}
	ratings = append(ratings, rating)
	if len(ratings) > config.Config.NodeRating.StorePoints {
		ratings = ratings[1:]
	}
	nodeRatings[network][index] = ratings

	logrus.Infof("    %s.%d ratings: %v", network, index, ratings)
}

func GetRatings(network string, index int) []int64 {
	logrus.Infof("GetRatings(%s.%d) RLock", network, index)
	mu.RLock()
	defer mu.RUnlock()

	logrus.Infof("GetRatings(%s.%d) RLocked", network, index)
	return nodeRatings[network][index]
}

func GetRating(network string, index int) int64 {
	logrus.Infof("GetRating(%s.%d) RLock", network, index)

	mu.RLock()
	defer mu.RUnlock()

	logrus.Infof("GetRating(%s.%d) RLocked", network, index)

	var rating int64
	for _, ra := range nodeRatings[network][index] {
		rating += ra
	}
	if len(nodeRatings[network][index]) == 0 {
		// downgrade nodes without any points
		return config.Config.NodeRating.ErrorRating
	}
	return rating
}

// returns N top nodes for specific network
func RoutableNodesWithBestRatings(network string) []NodeWithRating {
	networkConfig := config.Config.NetworksConfig[network]

	sortedNodes := NodesSortedByRating(network)
	// filter routable nodes only
	sortedRoutableNodes := []NodeWithRating{}
	for _, node := range sortedNodes {
		if networkConfig.Nodes[node.Index].AllowRouting {
			sortedRoutableNodes = append(sortedRoutableNodes, node)
		}
	}

	topNodes := []NodeWithRating{}
	for pos, node := range sortedRoutableNodes {
		// allow no more than MAX nodes in list
		if pos >= networkConfig.Rules.RoutingNodesMax {
			break
		}
		// include first MIN nodes anyway
		if pos < networkConfig.Rules.RoutingNodesMin {
			topNodes = append(topNodes, node)
			continue
		}
		if node.Rating <= networkConfig.Rules.GoodNodeMaxRating {
			topNodes = append(topNodes, node)
		}
	}

	return topNodes
}

func NodesSortedByRating(network string) []NodeWithRating {
	mu.RLock()
	defer mu.RUnlock()

	nodesWithRating := []NodeWithRating{}
	for index := range nodeRatings[network] {
		nodesWithRating = append(nodesWithRating, NodeWithRating{Index: index, Rating: GetRating(network, index), Network: network})
	}

	sort.SliceStable(nodesWithRating, func(i, j int) bool {
		return nodesWithRating[i].Rating < nodesWithRating[j].Rating
	})
	return nodesWithRating
}
