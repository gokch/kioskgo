package api

import (
	"net/http"

	"github.com/gokch/ipfs_mount/rpc/rpcconnect"
)

func RegisterAPI(mux *http.ServeMux) {

	filePath, fileHandler := rpcconnect.NewFileServiceHandler(NewFileServiceApi("/"))
	mux.Handle(filePath, fileHandler)

}
