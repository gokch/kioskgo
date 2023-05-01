package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gokch/kioskgo/p2p"
	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog"
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

	logger     = &zerolog.Logger{}
	rootPath   string
	privateKey string
	timeout    int64

	peerIds []string
	cids    []string
	paths   []string
)

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&rootPath, "rootpath", "r", "./", "root path")
	fs.StringVarP(&privateKey, "key", "k", "", "private key")
	fs.Int64VarP(&timeout, "timeout", "t", 0, "timeout seconds, 0 is no timeout")
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
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var ctx context.Context
	var cancel context.CancelFunc

	ctx = context.Background()
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	}
	if cancel != nil {
		defer cancel()
	}

	logger.Info().Msg("start client")

	client, err := p2p.NewClient(context.Background(), &p2p.ClientConfig{
		RootPath:   rootPath,
		Peers:      peerIds,
		PrivateKey: privateKey,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("init clinet is failed")
		return
	}
	defer client.Close()

	for i := range cids {
		path := paths[i]
		ci := cid.MustParse(cids[i])
		err = client.ReqDownload(ctx, ci, path)
		if err != nil {
			logger.Warn().Err(err).Str("cid", cids[i]).Msg("get cld is failed")
			return
		}
	}

	logger.Info().Msg("download is all done")

	handleKillSig(func() {
		client.Close()
	}, &logger)

	for {
		time.Sleep(time.Second)
	}

}

func handleKillSig(handler func(), logger *zerolog.Logger) {
	sigChannel := make(chan os.Signal, 1)

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for signal := range sigChannel {
			logger.Info().Msgf("Receive signal %s, Shutting down...", signal)
			handler()
			os.Exit(1)
		}
	}()
}
