package p2p

import (
	"context"
	"io"
	"time"

	"github.com/gokch/ipfs_mount/file"
	"github.com/gokch/ipfs_mount/mount"
	"github.com/ipfs/boxo/bitswap"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	dsync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-log"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/panjf2000/ants"
)

var (
	logger = log.Logger("client")
)

type ClientConfig struct {
	RootPath   string
	Peers      []string
	SizeWorker int
	ExpireSec  int
}

// waitlist 발신 ( 수신 / 송신 )
type Client struct {
	dag *mount.Dag
	mq  *ants.Pool

	host  host.Host
	bswap *bitswap.Bitswap
}

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
		dag: dag,
		mq:  mq,

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

func (c *Client) Self() string {
	return getHostAddress(c.host)
}

func (c *Client) Connect(ctx context.Context, targetPeer string) error {
	addrInfo, err := encodeAddrInfo(targetPeer)
	if err != nil {
		return err
	}

	return c.host.Connect(ctx, *addrInfo)
}

func (c *Client) Disconnect(ctx context.Context, targetPeer string) error {
	addrInfo, err := encodeAddrInfo(targetPeer)
	if err != nil {
		return err
	}
	return c.host.Network().ClosePeer(addrInfo.ID)
}

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

func (c *Client) Upload(ctx context.Context, path string, reader io.Reader) error {
	return c.mq.Submit(func() {
		cid, err := c.dag.Upload(ctx, path, reader)
		if err != nil {
			logger.Errorf("upload is failed | path : ", path, " | err : ", err.Error())
		} else {
			logger.Infof("upload is succeed | path : ", path, " | cid : ", cid.String())
		}
	})
}
