package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"
	"github.com/gokch/kioskgo/mount"
	"github.com/ipfs/boxo/bitswap"
	"github.com/ipfs/boxo/bitswap/client"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// waitlist 발신 ( 수신 / 송신 )
type Client struct {
	mount  *mount.Mount
	Client *client.Client

	host  host.Host
	bswap *bitswap.Bitswap

	havelist *file.FileManager // 현재 보유 목록
	waitlist *file.FileManager // 다운로드 대기 목록
}

func NewClient(ctx context.Context, address string, rootPath string) (*Client, error) {
	fs := file.NewFileStore(rootPath)
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(dsync.MutexWrap(datastore.NewMapDatastore())))

	host, err := makeHost(address, 0)
	if err != nil {
		return nil, err
	}
	address = getHostAddress(host)

	dht, err := dual.New(ctx, host)
	if err != nil {
		return nil, err
	}
	bsn := bsnet.NewFromIpfsHost(host, dht)
	bswap := bitswap.New(ctx, bsn, bs)

	// init bitswap
	mount, err := mount.NewMount(ctx, fs, bs, bswap)
	if err != nil {
		return nil, err
	}

	// init waitlist, havelist
	waitlist := file.NewFileManager(rootPath)
	havelist := file.NewFileManager(rootPath)

	// start server
	bsn.Start(bswap)

	return &Client{
		waitlist: waitlist,
		havelist: havelist,
		mount:    mount,
		host:     host,
		bswap:    bswap,
		Client:   bswap.Client,
	}, nil
}

func (c *Client) Connect(ctx context.Context, targetPeer string) error {
	maddr, err := multiaddr.NewMultiaddr(targetPeer)
	if err != nil {
		return err
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	// Directly connect to the peer that we know has the content
	// Generally this peer will come from whatever content routing system is provided, however go-bitswap will also
	// ask peers it is connected to for content so this will work
	if err := c.host.Connect(ctx, *info); err != nil {
		return err
	}

	return nil
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
