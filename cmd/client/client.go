package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gokch/kioskgo/p2p"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "client",
		Run: rootRun,
	}

	defaultConf = p2p.ClientConfig{
		RootPath: "./",
		Peers:    []string{},
	}

	rootPath   string
	timeout    int64
	workerSize int64
	expireSec  int64

	peerIds []string
	cids    []string
	paths   []string
)

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&rootPath, "rootpath", "r", "./", "root path")
	fs.Int64VarP(&timeout, "timeout", "t", 0, "timeout seconds, 0 is no timeout")
	fs.Int64VarP(&workerSize, "worker", "w", 1, "worker size")
	fs.Int64VarP(&expireSec, "expire", "e", 600, "expire seconds")

	fs.StringArrayVar(&peerIds, "peers", []string{}, "connect peer id")
	fs.StringArrayVarP(&cids, "cids", "c", []string{}, "download cid")
	fs.StringArrayVarP(&paths, "paths", "p", []string{}, "download path per cid")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	log.SetLogLevel("*", "debug")
	logger := log.Logger("client")

	var ctx context.Context
	var cancel context.CancelFunc

	ctx = context.Background()
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	}
	if cancel != nil {
		defer cancel()
	}

	logger.Info("start client")

	client, err := p2p.NewClient(context.Background(), &p2p.ClientConfig{
		RootPath:   rootPath,
		Peers:      peerIds,
		SizeWorker: int(workerSize),
		ExpireSec:  int(expireSec),
	})
	if err != nil {
		logger.Warn(err, "init client is failed")
		return
	}

	for i := range cids {
		path := paths[i]
		ci, err := cid.Parse(cids[i])
		if err != nil {
			logger.Warn(err, fmt.Sprint("cid :", cids[i], "invalid cid"))
			return
		}
		err = client.Download(ctx, ci, path)
		if err != nil {
			logger.Warn(err, fmt.Sprint("cid :", cids[i], "get cid is failed"))
			return
		}
	}

	handleKillSig(func() {
		client.Close()
	}, logger)

	client.Close()
	logger.Info("download is all done")
}

func handleKillSig(handler func(), logger log.EventLogger) {
	sigChannel := make(chan os.Signal, 1)

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for signal := range sigChannel {
			logger.Info("Receive signal %s, Shutting down...", signal)
			handler()
			os.Exit(1)
		}
	}()
}
