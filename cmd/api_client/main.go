package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/gokch/ipfs_mount/rpc"
	"github.com/gokch/ipfs_mount/rpc/rpcconnect"
)

type Client struct {
	client rpcconnect.ClientServiceClient
}

func NewClient(addr string) *Client {
	return &Client{
		client: rpcconnect.NewClientServiceClient(http.DefaultClient, addr),
	}
}

func main() {
	ctx := context.Background()

	client := NewClient("http://localhost:9876")
	resp, err := client.client.Connect(ctx, connect.NewRequest(&rpc.ConnectRequest{
		Peers: []*rpc.Peer{
			{
				Peerid: "",
			},
		},
	}))
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	fmt.Println(resp.Msg)

	resp2, err := client.client.Download(ctx, connect.NewRequest(&rpc.DownloadRequest{
		Files: []*rpc.File{
			{
				Cid:  "",
				Path: "/path/angmond.jpg",
			},
		},
	}))
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	fmt.Println(resp2.Msg)
	time.Sleep(time.Second * 100)

}
