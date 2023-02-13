package node_rating

import (
	"fmt"
	"node-balancer/internal/server/config"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	nodeRatings = map[string][]int64{}
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

func ratingKey(network string, index int) string {
	return fmt.Sprintf("%s.%d", network, index)
}

func AddRating(network string, index int, blockDelayBlocks int64, err error) {
	rating := blockDelayBlocks
	if err != nil {
		rating = config.Config.NodeRating.ErrorRating
	}

	mu.Lock()
	defer mu.Unlock()

	ratings := nodeRatings[ratingKey(network, index)]

	ratings = append(ratings, rating)
	if len(ratings) > config.Config.NodeRating.StorePoints {
		ratings = ratings[1:]
	}
	nodeRatings[ratingKey(network, index)] = ratings

	logrus.Infof("    %s.%d ratings: %v", network, index, ratings)
}

func getRatings(network string, index int) []int64 {
	mu.RLock()
	defer mu.RUnlock()

	return nodeRatings[ratingKey(network, index)]
}

func getRating(network string, index int) int64 {
	mu.RLock()
	defer mu.RUnlock()

	ratings := nodeRatings[ratingKey(network, index)]
	if len(ratings) == 0 {
		// downgrade nodes without any points
		return config.Config.NodeRating.ErrorRating
	}
	var rating int64
	for _, ra := range ratings {
		rating += ra
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
	nodesWithRating := []NodeWithRating{}

	for index := range config.Config.NetworksConfig[network].Nodes {
		rating := getRating(network, index)
		nodesWithRating = append(nodesWithRating, NodeWithRating{Index: index, Rating: rating, Network: network})
	}

	sort.SliceStable(nodesWithRating, func(i, j int) bool {
		return nodesWithRating[i].Rating < nodesWithRating[j].Rating
	})
	return nodesWithRating
}
