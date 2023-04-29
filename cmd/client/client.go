package main

import (
	"context"
	"fmt"
	"os"

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
	client, err := p2p.NewClient(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	err = client.Connect(cmd.Context(), "/ip4/172.30.1.254/tcp/54264/p2p/QmQQboMmhTBVfNBvztEbsC8VZxeXvb11SaTftJ6EaYcm14")
	if err != nil {
		fmt.Println(err)
		return
	}

	client.AddWaitlist(cmd.Context(), cid.MustParse("bafkreicsqaff7pryibb4lucdonapngzvk44nspdaoal3qn3oq55efix7kq"), "/tmp/test.txt")

}
