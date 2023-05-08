package api

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/gokch/ipfs_mount/file"
	"github.com/gokch/ipfs_mount/rpc"
	"github.com/gokch/ipfs_mount/rpc/rpcconnect"
)

var _ rpcconnect.FileServiceClient = (*FileServiceApi)(nil)

func NewFileServiceApi(rootPath string) *FileServiceApi {
	return &FileServiceApi{
		fs: file.NewFileStore(rootPath),
	}
}

type FileServiceApi struct {
	fs *file.FileStore
}

func (f *FileServiceApi) Upload(ctx context.Context, req *connect.Request[rpc.UploadRequest]) (*connect.Response[rpc.UploadResponse], error) {
	return nil, nil
}

func (f *FileServiceApi) Download(context.Context, *connect.Request[rpc.DownloadRequest]) (*connect.Response[rpc.DownloadResponse], error) {

	return nil, nil
}
