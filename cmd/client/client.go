package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gokch/kioskgo/file"
	"github.com/gokch/kioskgo/p2p"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "client",
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
	fs := file.NewFileStore("client")
	client, err := p2p.NewP2PClient(context.Background(), "", fs)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	err = client.Connect(cmd.Context(), "/ip4/127.0.0.1/tcp/1607/p2p/12D3KooWHq35S2aGStH1kG4LH99TEi8SBLUvbanVJADRoBUNPDcP")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Download(cmd.Context(), cid.MustParse("bafkreicsqaff7pryibb4lucdonapngzvk44nspdaoal3qn3oq55efix7kq"), "/tmp/test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

}
