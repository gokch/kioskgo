package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gokch/kioskgo/file"
	"github.com/gokch/kioskgo/p2p"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "server",
		Run: rootRun,
	}
)

func init() {
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	fs := file.NewFileStore("server")
	server, err := p2p.NewP2P(cmd.Context(), "", fs, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	cid, err := server.Upload(cmd.Context(), file.NewReaderFromBytes([]byte("testaaaaaa")))
	fmt.Println(cid.String())
	fmt.Println(server.Address)

	handleKillSig(func() {
		server.Close()
	})

	for {
		time.Sleep(time.Second)
	}

}

func handleKillSig(handler func()) {
	sigChannel := make(chan os.Signal, 1)

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for signal := range sigChannel {
			fmt.Println(signal)
			handler()
			os.Exit(1)
		}
	}()
}
