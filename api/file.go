package api

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/gokch/kioskgo/file"
	"github.com/gokch/kioskgo/types"
	"github.com/gokch/kioskgo/types/typesconnect"
)

var _ typesconnect.FileServiceClient = (*FileServiceApi)(nil)

func NewFileServiceApi(rootPath string) *FileServiceApi {
	return &FileServiceApi{
		fs: file.NewFileSystem(rootPath),
	}
}

type FileServiceApi struct {
	fs *file.FileSystem
}

func (f *FileServiceApi) Upload(ctx context.Context, req *connect.Request[types.UploadReq]) (*connect.Response[types.UploadRes], error) {
	return nil, nil
}

func (f *FileServiceApi) Download(context.Context, *connect.Request[types.DownloadReq]) (*connect.Response[types.DownloadRes], error) {

	return nil, nil
}
