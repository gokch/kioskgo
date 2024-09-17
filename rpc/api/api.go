package api

import (
	"net/http"

	"github.com/rabbitprincess/ipfs_mount/p2p"
	"github.com/rabbitprincess/ipfs_mount/rpc"
	"github.com/rabbitprincess/ipfs_mount/rpc/rpcconnect"
)

func RegisterAPI(mux *http.ServeMux, client *p2p.Client) {
	clientPath, clientHandler := rpcconnect.NewClientServiceHandler(NewClientServiceApi(client))
	mux.Handle(clientPath, clientHandler)

	if client.Role == rpc.Role_SERVER {

	}
}
