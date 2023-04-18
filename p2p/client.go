package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"
	"github.com/ipfs/boxo/routing/http/contentrouter"
	"github.com/ipfs/go-cid"
)

// waitlist 발신 ( 수신 / 송신 )
type Client struct {
	*P2P
	havelist *file.FileManager // 현재 보유 목록
	waitlist *file.FileManager // 다운로드 대기 목록
}

func NewClient(ctx context.Context, address string, rootPath string, clientrouter contentrouter.Client) (*Client, error) {
	if clientrouter == nil {
	}

	p2p, err := NewP2P(ctx, address, rootPath, clientrouter)
	if err != nil {
		return nil, err
	}
	// init waitlist, havelist
	waitlist := file.NewFileManager(rootPath)
	havelist := file.NewFileManager(rootPath)
	readers, err := p2p.fs.Iterate("")
	if err != nil {
		return nil, err
	}
	for _, reader := range readers {
		havelist.PutReader(reader)
	}

	return &Client{
		waitlist: waitlist,
		havelist: havelist,
		P2P:      p2p,
	}, nil
}

func (c *Client) Start() {

}

// 1. 클라이언트가 특정 피어를 가지고 싶다고 요청
func (c *Client) AddWaitlist(cid cid.Cid, path string) {
	// c.waitlist.Add(cid, path)
}

// 1. waitlist 에서 요청한 cid 를 다운로드 받았을 경우
// 2. 새로운 cid 를 가진 파일을 피어에 올릴 경우
func (c *Client) AddHavelist(cid cid.Cid, path string) {
	// c.havelist.Add(cid, path)
}

func (c *Client) RecvDownload(ctx context.Context, cid cid.Cid, path string) error {
	err := c.P2P.Download(ctx, cid, path)
	if err != nil {
		return err
	}

	// 다운로드가 끝났을 시 waitlist 에서 지운다 + havelist 에 추가한다
	// c.waitlist.Remove(cid, path)
	c.AddHavelist(cid, path)

	return nil
}

// peer 에 제공? - TODO: IPNS 를 통해 고유값 지정이 필요!!
func (c *Client) RecvUpload(ctx context.Context, path string) {
	c.P2P.Upload(ctx, path)
}
