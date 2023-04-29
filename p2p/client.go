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

// waitlist 발신 ( 수신 / 송신 )
type Client struct {
	mount    *mount.Mount
	havelist *file.FileManager // 현재 보유 목록
	waitlist *file.FileManager // 다운로드 대기 목록

	host  host.Host
	bswap *bitswap.Bitswap
}

func NewClient(ctx context.Context, rootPath string) (*Client, error) {
	// init fs
	fs := file.NewFileStore(rootPath)

	// init memory bs for dht
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(dsync.MutexWrap(datastore.NewMapDatastore())))
	host, err := makeHost(rootPath)
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

	// init waitlist, havelist
	waitlist := file.NewFileManager(rootPath)
	havelist := file.NewFileManager(rootPath)

	return &Client{
		waitlist: waitlist,
		havelist: havelist,
		mount:    mount,
		host:     host,
		bswap:    bswap,
	}, nil
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

// 1. waitlist 에서 요청한 cid 를 다운로드 받았을 경우
// 2. 새로운 cid 를 가진 파일을 피어에 올릴 경우
func (c *Client) AddHavelist(ctx context.Context, cid cid.Cid, path string) {
	c.havelist.Put(path, cid)
}

func (c *Client) RecvDownload(ctx context.Context, cid cid.Cid, path string) error {
	err := c.mount.Download(ctx, cid, path)
	if err != nil {
		return err
	}

	// 다운로드가 끝났을 시 waitlist 에서 지운다 + havelist 에 추가한다
	c.waitlist.Delete(path, cid)
	c.AddHavelist(ctx, cid, path)

	return nil
}

// peer 에 제공? - TODO: IPNS 를 통해 고유값 지정이 필요!!
func (c *Client) RecvUpload(ctx context.Context, path string) {
	c.mount.Upload(ctx, path)
}
