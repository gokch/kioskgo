package api

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/gokch/ipfs_mount/p2p"
	"github.com/gokch/ipfs_mount/rpc"
	"github.com/gokch/ipfs_mount/rpc/rpcconnect"
	"github.com/ipfs/go-cid"
)

var _ rpcconnect.ClientServiceClient = (*ClientServiceApi)(nil)

func NewClientServiceApi(client *p2p.Client) *ClientServiceApi {
	return &ClientServiceApi{
		client: client,
	}
}

type ClientServiceApi struct {
	client *p2p.Client
}

func (c *ClientServiceApi) Connect(ctx context.Context, req *connect.Request[rpc.ConnectRequest]) (*connect.Response[rpc.ConnectResponse], error) {
	conn := make([]*rpc.Peer, 0, 10)
	unconn := make([]*rpc.Peer, 0, 10)

	for _, peer := range req.Msg.GetPeers() {
		if err := c.client.Connect(ctx, peer.GetPeerid()); err != nil {
			unconn = append(unconn, peer)
		} else {
			conn = append(conn, peer)
		}
	}
	conResp := &rpc.ConnectResponse{
		Response: &rpc.Response{},
		Succeed:  conn,
		Failed:   unconn,
	}

	return connect.NewResponse(conResp), nil
}

func (c *ClientServiceApi) Disconnect(context.Context, *connect.Request[rpc.DisconnectRequest]) (*connect.Response[rpc.DisconnectResponse], error) {
	return nil, nil
}

func (c *ClientServiceApi) IsConnect(context.Context, *connect.Request[rpc.IsConnectRequest]) (*connect.Response[rpc.IsConnectResponse], error) {
	return nil, nil
}

func (c *ClientServiceApi) Upload(ctx context.Context, req *connect.Request[rpc.UploadRequest]) (*connect.Response[rpc.UploadResponse], error) {
	succeed := make([]*rpc.File, 0, 10)
	failed := make([]*rpc.File, 0, 10)

	for _, file := range req.Msg.GetFiles() {
		cid, err := cid.Parse(file.GetCid())
		if err != nil {
			failed = append(failed, file)
			continue
		}

		err = c.client.Upload(ctx, cid, file.GetPath())
		if err != nil {
			failed = append(failed, file)
			continue
		}
		succeed = append(succeed, file)

	}
	uploadResp := &rpc.UploadResponse{
		Response: &rpc.Response{},
		Succeed:  succeed,
		Failed:   failed,
	}

	return connect.NewResponse(uploadResp), nil
}

func (c *ClientServiceApi) Download(ctx context.Context, req *connect.Request[rpc.DownloadRequest]) (*connect.Response[rpc.DownloadResponse], error) {
	succeed := make([]*rpc.File, 0, 10)
	failed := make([]*rpc.File, 0, 10)

	for _, file := range req.Msg.GetFiles() {
		cid, err := cid.Parse(file.GetCid())
		if err != nil {
			failed = append(failed, file)
			continue
		}

		err = c.client.Download(ctx, cid, file.GetPath())
		if err != nil {
			failed = append(failed, file)
			continue
		}
		succeed = append(succeed, file)

	}
	downloadResp := &rpc.DownloadResponse{
		Response: &rpc.Response{},
		Succeed:  succeed,
		Failed:   failed,
	}

	return connect.NewResponse(downloadResp), nil
}
