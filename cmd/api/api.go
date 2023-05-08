package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gokch/ipfs_mount/p2p"
	"github.com/gokch/ipfs_mount/rpc/api"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "api",
		Run: rootRun,
	}

	addr       string
	rootPath   string
	timeout    int64
	workerSize int64
	expireSec  int64

	peerIds []string
)

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&addr, "addr", "a", "localhost:9876", "address")
	fs.StringVarP(&rootPath, "rootpath", "r", "./", "root path")
	fs.Int64VarP(&timeout, "timeout", "t", 0, "timeout seconds, 0 is no timeout")
	fs.Int64VarP(&workerSize, "worker", "w", 1, "worker size")
	fs.Int64VarP(&expireSec, "expire", "e", 600, "expire seconds")

	fs.StringArrayVar(&peerIds, "peers", []string{}, "connect peer id")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	client, err := p2p.NewClient(cmd.Context(), &p2p.ClientConfig{
		RootPath:   rootPath,
		Peers:      peerIds,
		SizeWorker: int(workerSize),
		ExpireSec:  int(expireSec),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	mux := http.NewServeMux()
	api.RegisterAPI(mux, client)

	go func() {
		err = http.ListenAndServe(addr, mux)
		if err != nil {
			fmt.Println(err)
		}
	}()

	i := handleKillSig(func() {
		client.Close()
	})
	<-i.C
}

type interrupt struct {
	C chan struct{}
}

func handleKillSig(handler func()) interrupt {
	i := interrupt{
		C: make(chan struct{}),
	}
	sigChannel := make(chan os.Signal, 1)

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for signal := range sigChannel {
			fmt.Println(signal)
			handler()
			os.Exit(1)
		}
	}()
	return i
}
