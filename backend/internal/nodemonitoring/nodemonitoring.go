package nodemonitoring

import (
	"context"
	"crypto-exporter/internal/metrics"
	"crypto-exporter/internal/server/config"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"math/big"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

func Run(config config.Config) {

	//logrus.Printf("config: %#v\n", config.RpcConfig)

	for range time.Tick(time.Millisecond * 2000) {
		for _, s := range config.RpcConfig {
			for _, v := range s {
				go printBlockNumber(v)
			}
		}
		fmt.Println("\n---")
	}

}

func printBlockNumber(s config.Node) {
	start := time.Now()
	_, blockNum := getBlockNumber(s)
	spent := time.Since(start)

	//err, syncProgress := getSyncStatus(s)
	//logrus.Println(syncProgress)

	err, getLastBlock, getLastBlockTime := getLastKnowBlock(s)
	if err != nil {

	} else {
		metrics.OpsBlockHight.WithLabelValues("polygon", s.Name, s.Url, strconv.FormatBool(s.Public)).Set(float64(getLastBlock.Int64()))
	}

	//labels := make(map[string]string)
	//labels["name"] = s.Name
	//v, _ := metrics.OpsBlockHight.GetMetricWith(labels)
	//logrus.Error(v)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 15, 20, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintf(w, "\n %v\t %v\t %v\t %v\t %v\t %s ", start.Format("15:04:05.99999"), spent, blockNum, getLastBlock, getLastBlockTime.Format("15:04:05.99999"), s.Name)
}

func getBlockNumber(s config.Node) (error, uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	client, err := ethclient.DialContext(ctx, s.Url)
	if err != nil {
		logrus.Errorf("\n Error connect to: %v Error: %v", s.Url, err)
	}

	defer cancel()
	header, err := client.BlockNumber(ctx)
	if err != nil {
		logrus.Errorf("\nGet block num. Error: %v", err)
		return err, 0
	} else {
		return nil, header
	}
}

func getLastKnowBlock(s config.Node) (error, *big.Int, time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	client, err := ethclient.DialContext(ctx, s.Url)
	if err != nil {
		logrus.Errorf("\nError connect to: %v Error: %v", s.Url, err)
	}

	defer cancel()
	latesHeader, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		logrus.Errorf("\n Get last know block num. %v", err)
		return err, nil, time.Unix(0, 0)
	} else {
		blockTime := time.Unix(int64(latesHeader.Time), 0)
		return err, latesHeader.Number, blockTime
	}
}

func getSyncStatus(s config.Node) (error, uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	client, err := ethclient.DialContext(ctx, s.Url)
	if err != nil {
		logrus.Errorf("\nError connect to: %v Error: %v", s.Url, err)
	}

	defer cancel()
	SyncProgress, err := client.SyncProgress(ctx)
	if err != nil {
		logrus.Errorf("\n Get Sync progress. %v", err)
		return err, 0
	} else {
		//SyncStatus := SyncProgress.
		logrus.Println(SyncProgress)
		return err, 0
	}
}
