package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/gokch/kioskgo/file"
	"github.com/gokch/kioskgo/mount"
	"github.com/ipfs/boxo/bitswap"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	"github.com/ipfs/go-cid"
	dsync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/panjf2000/ants"
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
	cs := file.NewCacheStore()

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

func (c *Client) Download(ctx context.Context, cid cid.Cid, path string) error {
	return c.mq.Submit(func() {
		err := c.dag.Download(ctx, cid, path)
		if err != nil {
			fmt.Println("download not finished")
		} else {
			fmt.Println("download finished")
		}
	})
}

func (c *Client) Upload(ctx context.Context, cid cid.Cid, path string) error {
	return c.mq.Submit(func() {
		cid, err := c.dag.Upload(ctx, path, nil)
		if err != nil {
			fmt.Println("upload not finished")
		} else {
			fmt.Println("upload finished | cid : ", cid.String())
		}
	})
}
