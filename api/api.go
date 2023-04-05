package api

import (
	"net/http"

	"github.com/gokch/kioskgo/types/typesconnect"
)

func RegisterAPI(mux *http.ServeMux) {

	filePath, fileHandler := typesconnect.NewFileServiceHandler(NewFileServiceApi("/"))
	mux.Handle(filePath, fileHandler)

}
