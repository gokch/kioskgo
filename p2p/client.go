package p2p

import (
	"context"
	"time"

	"github.com/ipfs/boxo/bitswap"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	dsync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-log"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/panjf2000/ants"
	"github.com/rabbitprincess/ipfs_mount/file"
	"github.com/rabbitprincess/ipfs_mount/mount"
	"github.com/rabbitprincess/ipfs_mount/rpc"
)

var (
	logger = log.Logger("client")
)

// ClientConfig contains the configuration for the client.
type ClientConfig struct {
	// RootPath is the path to the root directory of the file system.
	RootPath string
	// Peers is a list of peer addresses to connect to.
	Peers []string
	// SizeWorker is the number of worker goroutines to use for downloading and uploading files.
	SizeWorker int
	// ExpireSec is the number of seconds to keep a file in the cache before it is evicted.
	ExpireSec int

	Role rpc.Role
}

// Client is a peer-to-peer client that can connect to other peers in the network, download and upload files, and manage its own connections.
type Client struct {
	Role rpc.Role

	// dag is the dag mount.
	dag *mount.Dag
	// mq is the ants pool.
	mq *ants.Pool

	// host is the libp2p host.
	host host.Host
	// bswap is the bitswap client.
	bswap *bitswap.Bitswap
}

// NewClient creates a new client.
func NewClient(ctx context.Context, cfg *ClientConfig) (*Client, error) {
	// init fs
	fs := file.NewFileStore(cfg.RootPath)
	fm := file.NewFileManager()
	cs := file.NewCacheStore(time.Second * time.Duration(cfg.ExpireSec))

	// init memory bs for dht
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(dsync.MutexWrap(cs)))
	host, err := makeHost(cfg.RootPath)
	if err != nil {
		return nil, err
	}

	// init dht for bitswap network
	ipfsdht, err := dht.New(ctx, host)
	if err != nil {
		return nil, err
	}

	// init bitswap network
	bsn := bsnet.NewFromIpfsHost(host, ipfsdht)
	bswap := bitswap.New(ctx, bsn, bs)
	bsn.Start(bswap)

	// init dag
	dag, err := mount.NewDag(ctx, -1, mount.NewMount(fs, fm, bs), bswap)
	if err != nil {
		return nil, err
	}

	mq, err := ants.NewPool(cfg.SizeWorker, ants.WithExpiryDuration(time.Second*time.Duration(cfg.ExpireSec)))
	if err != nil {
		return nil, err
	}

	c := &Client{
		Role: cfg.Role,
		dag:  dag,
		mq:   mq,

		host:  host,
		bswap: bswap,
	}

	// connect
	for _, peer := range cfg.Peers {
		err := c.Connect(ctx, peer)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Self returns the client's self address.
func (c *Client) Self() string {
	return getHostAddress(c.host)
}

// Connect connects the client to a peer.
func (c *Client) Connect(ctx context.Context, targetPeer string) error {
	addrInfo, err := encodeAddrInfo(targetPeer)
	if err != nil {
		return err
	}

	return c.host.Connect(ctx, *addrInfo)
}

// Disconnect disconnects the client from a peer.
func (c *Client) Disconnect(ctx context.Context, targetPeer string) error {
	addrInfo, err := encodeAddrInfo(targetPeer)
	if err != nil {
		return err
	}
	return c.host.Network().ClosePeer(addrInfo.ID)
}

// Connect connects the client to a peer.
func (c *Client) IsConnect(ctx context.Context) map[string]bool {
	conns := make(map[string]bool)
	for _, conn := range c.host.Network().Conns() {
		conns[conn.ID()] = true
	}
	return conns
}

// Close closes the client.
func (c *Client) Close() error {
	for c.mq.Running() > 0 {
		time.Sleep(time.Second)
	}

	c.mq.Release()

	if err := c.bswap.Close(); err != nil {
		return err
	}
	if err := c.host.Close(); err != nil {
		return err
	}
	return nil
}

// Download downloads a file from the network.
func (c *Client) Download(ctx context.Context, cid string, path string) error {
	return c.mq.Submit(func() {
		err := c.dag.Download(ctx, cid, path)
		if err != nil {
			logger.Errorf("download is failed | cid : %s | path : %s | err : %s", cid, path, err.Error())
		} else {
			logger.Infof("download is succeed | cid : %s | path : %s", cid, path)
		}
	})
}

// Upload uploads a file to the network.
func (c *Client) Upload(ctx context.Context, path string, reader *file.Reader) error {
	return c.mq.Submit(func() {
		cid, err := c.dag.Upload(ctx, path, reader)
		if err != nil {
			logger.Errorf("upload is failed | path : ", path, " | err : ", err.Error())
		} else {
			logger.Infof("upload is succeed | path : %s | cid : %s", path, cid.String())
		}
	})
}
