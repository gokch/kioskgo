package api

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/rabbitprincess/ipfs_mount/p2p"
	"github.com/rabbitprincess/ipfs_mount/rpc"
	"github.com/rabbitprincess/ipfs_mount/rpc/rpcconnect"
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

func (c *ClientServiceApi) Disconnect(ctx context.Context, req *connect.Request[rpc.DisconnectRequest]) (*connect.Response[rpc.DisconnectResponse], error) {
	disconn := make([]*rpc.Peer, 0, 10)
	undisconn := make([]*rpc.Peer, 0, 10)

	for _, peer := range req.Msg.Peers {
		if err := c.client.Disconnect(ctx, peer.GetPeerid()); err != nil {
			undisconn = append(undisconn, peer)
		} else {
			disconn = append(disconn, peer)
		}
	}
	conResp := &rpc.DisconnectResponse{
		Response: &rpc.Response{},
		Succeed:  disconn,
		Failed:   undisconn,
	}

	return connect.NewResponse(conResp), nil
}

func (c *ClientServiceApi) IsConnect(ctx context.Context, req *connect.Request[rpc.IsConnectRequest]) (*connect.Response[rpc.IsConnectResponse], error) {
	conn := make([]*rpc.Peer, 0, 10)
	unconn := make([]*rpc.Peer, 0, 10)

	conns := c.client.IsConnect(ctx)
	for _, con := range req.Msg.Peers {
		if _, ok := conns[con.GetPeerid()]; ok {
			conn = append(conn, con)
		} else {
			unconn = append(unconn, con)
		}
	}

	conResp := &rpc.IsConnectResponse{
		Response:   &rpc.Response{},
		Connects:   conn,
		Unconnects: unconn,
	}

	return connect.NewResponse(conResp), nil
}

// TODO
func (c *ClientServiceApi) Query(ctx context.Context, req *connect.Request[rpc.QueryRequest]) (*connect.Response[rpc.QueryResponse], error) {

	queryResp := &rpc.QueryResponse{}
	return connect.NewResponse(queryResp), nil
}

func (c *ClientServiceApi) Upload(ctx context.Context, req *connect.Request[rpc.UploadRequest]) (*connect.Response[rpc.UploadResponse], error) {
	succeed := make([]*rpc.File, 0, 10)
	failed := make([]*rpc.File, 0, 10)

	for _, file := range req.Msg.GetFiles() {
		err := c.client.Upload(ctx, file.GetPath(), nil)
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
		err := c.client.Download(ctx, file.GetCid(), file.GetPath())
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
