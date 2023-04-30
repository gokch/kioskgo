package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"
	"github.com/gokch/kioskgo/mount"
	"github.com/ipfs/boxo/bitswap"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

type ClientConfig struct {
	RootPath   string
	Peers      []string
	PrivateKey string
}

// waitlist 발신 ( 수신 / 송신 )
type Client struct {
	mount    *mount.Mount
	waitlist *file.FileManager // 다운로드 대기 목록

	host  host.Host
	bswap *bitswap.Bitswap
}

func NewClient(ctx context.Context, cfg *ClientConfig) (*Client, error) {
	// init fs
	fs := file.NewFileStore(cfg.RootPath)

	// init waitlist
	waitlist := file.NewFileManager(cfg.RootPath)

	// init memory bs for dht
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(dsync.MutexWrap(datastore.NewMapDatastore())))
	host, err := makeHost(cfg.PrivateKey)
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

	// init mount
	mount, err := mount.NewMount(ctx, fs, bs, bswap)
	if err != nil {
		return nil, err
	}

	c := &Client{
		waitlist: waitlist,
		mount:    mount,
		host:     host,
		bswap:    bswap,
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
	if err := c.bswap.Close(); err != nil {
		return err
	}
	if err := c.host.Close(); err != nil {
		return err
	}
	return nil
}

// 1. 클라이언트가 특정 피어를 가지고 싶다고 요청
func (c *Client) AddWaitlist(ctx context.Context, cid cid.Cid, path string) {
	c.waitlist.Put(path, cid)
}

func (c *Client) RecvDownload(ctx context.Context, cid cid.Cid, path string) error {
	err := c.mount.Download(ctx, cid, path)
	if err != nil {
		return err
	}
	// TODO : wantlist 삭제, BlockReceivedNotifier 콜백 사용
	// c.bswap.Client.ReceiveMessage()

	// 다운로드가 끝났을 시 waitlist 에서 지운다
	c.waitlist.Delete(path, cid)
	return nil
}

// peer 에 제공?
func (c *Client) RecvUpload(ctx context.Context, path string) (cid.Cid, error) {
	ci, err := c.mount.Upload(ctx, path)
	if err != nil {
		return cid.Cid{}, err
	}
	return ci, nil
}
