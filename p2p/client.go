package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"
	"github.com/ipfs/boxo/bitswap/client"
	"github.com/ipfs/go-cid"
)

// waitlist 발신 ( 수신 / 송신 )
type Client struct {
	*Mount
	Client *client.Client

	havelist *file.FileManager // 현재 보유 목록
	waitlist *file.FileManager // 다운로드 대기 목록
}

func NewClient(ctx context.Context, address string, rootPath string) (*Client, error) {
	// TODO : init bitswap ( or offline. anyway. 빨리 고쳐라. )
	mount, err := NewMount(ctx, rootPath, nil)
	if err != nil {
		return nil, err
	}

	// init waitlist, havelist
	waitlist := file.NewFileManager(rootPath)
	havelist := file.NewFileManager(rootPath)

	return &Client{
		waitlist: waitlist,
		havelist: havelist,
		Mount:    mount,
		// Client:   p2p.bswap.Client,
	}, nil
}

// func (p *P2P) Connect(ctx context.Context, targetPeer string) error {
// 	maddr, err := multiaddr.NewMultiaddr(targetPeer)
// 	if err != nil {
// 		return err
// 	}
// 	info, err := peer.AddrInfoFromP2pAddr(maddr)
// 	if err != nil {
// 		return err
// 	}

// 	// Directly connect to the peer that we know has the content
// 	// Generally this peer will come from whatever content routing system is provided, however go-bitswap will also
// 	// ask peers it is connected to for content so this will work
// 	if err := p.host.Connect(ctx, *info); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (c *Client) Start() {

}

// 1. 클라이언트가 특정 피어를 가지고 싶다고 요청
func (c *Client) AddWaitlist(cid cid.Cid, path string) {
	c.waitlist.Put(path, cid)
}

// 1. waitlist 에서 요청한 cid 를 다운로드 받았을 경우
// 2. 새로운 cid 를 가진 파일을 피어에 올릴 경우
func (c *Client) AddHavelist(cid cid.Cid, path string) {
	c.havelist.Put(path, cid)
}

func (c *Client) RecvDownload(ctx context.Context, cid cid.Cid, path string) error {
	err := c.Mount.Download(ctx, cid, path)
	if err != nil {
		return err
	}

	// 다운로드가 끝났을 시 waitlist 에서 지운다 + havelist 에 추가한다
	c.waitlist.Delete(path, cid)
	c.AddHavelist(cid, path)

	return nil
}

// peer 에 제공? - TODO: IPNS 를 통해 고유값 지정이 필요!!
func (c *Client) RecvUpload(ctx context.Context, path string) {
	c.Mount.Upload(ctx, path)
}
