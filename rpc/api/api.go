package api

import (
	"net/http"

	"github.com/gokch/ipfs_mount/p2p"
	"github.com/gokch/ipfs_mount/rpc/rpcconnect"
)

func RegisterAPI(mux *http.ServeMux, client *p2p.Client) {
	clientPath, clientHandler := rpcconnect.NewClientServiceHandler(NewClientServiceApi(client))
	mux.Handle(clientPath, clientHandler)

	if client.IsServer {

	}
}
