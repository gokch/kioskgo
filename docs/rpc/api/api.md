# Package api

Package api provides a client service API for the ipfs_mount project.\
API methods are defined at [rpc/proto/client_service.proto](https://github.com/rabbitprincess/ipfs-mount/blob/main/rpc/proto/client_service.proto) also

## ClientServiceApi struct

The ClientServiceApi struct provides methods for connecting to, disconnecting from, and uploading and downloading files from ipfs nodes.

### Methods

* **Connect(ctx context.Context, req *connect.Request[rpc.ConnectRequest]) (*connect.Response[rpc.ConnectResponse], error)**

Connects to the specified peers.

* **Disconnect(context.Context, *connect.Request[rpc.DisconnectRequest]) (*connect.Response[rpc.DisconnectResponse], error)**

Disconnects from the specified peers.

* **IsConnect(context.Context, *connect.Request[rpc.IsConnectRequest]) (*connect.Response[rpc.IsConnectResponse], error)**

Checks if the client is connected to the specified peer.

* **Upload(ctx context.Context, req *connect.Request[rpc.UploadRequest]) (*connect.Response[rpc.UploadResponse], error)**

Uploads the specified files to the ipfs network.

* **Download(ctx context.Context, req *connect.Request[rpc.DownloadRequest]) (*connect.Response[rpc.DownloadResponse], error)**

Downloads the specified files from the ipfs network.

