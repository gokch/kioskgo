package p2p

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	client "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/go-libipfs/blocks"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multihash"
)

func NewIPFSstore(ctx context.Context, online bool) (*IPFSstore, error) {
	localApi, err := client.NewLocalApi()
	if err != nil {
		return nil, fmt.Errorf("getting ipfs api: %w", err)
	}
	api, err := localApi.WithOptions(options.Api.Offline(!online))
	if err != nil {
		return nil, fmt.Errorf("setting offline mode: %s", err)
	}
	ipfsStore := &IPFSstore{
		ctx: ctx,
		api: api,
	}
	return ipfsStore, nil
}

type IPFSstore struct {
	ctx context.Context
	api iface.CoreAPI
}

var _ blockstore.Blockstore = (*IPFSstore)(nil)

func (i *IPFSstore) Has(ctx context.Context, cid cid.Cid) (bool, error) {
	_, err := i.api.Block().Stat(ctx, path.IpldPath(cid))
	if err != nil {
		if err.Error() == "blockservice: key not found" {
			return false, nil
		}
		return false, fmt.Errorf("getting ipfs block: %w", err)
	}
	return true, nil
}

func (i *IPFSstore) Get(ctx context.Context, cid cid.Cid) (blocks.Block, error) {
	reader, err := i.api.Block().Get(ctx, path.IpldPath(cid))
	if err != nil {
		return nil, fmt.Errorf("getting ipfs block: %w", err)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	basicBlock, err := blocks.NewBlockWithCid(data, cid)
	if err != nil {
		return nil, err
	}
	return basicBlock, nil
}

func (i *IPFSstore) GetSize(ctx context.Context, cid cid.Cid) (int, error) {
	st, err := i.api.Block().Stat(ctx, path.IpldPath(cid))
	if err != nil {
		return 0, fmt.Errorf("getting ipfs block: %w", err)
	}

	return st.Size(), nil
}

func (i *IPFSstore) Put(ctx context.Context, block blocks.Block) error {
	mhd, err := multihash.Decode(block.Cid().Hash())
	if err != nil {
		return err
	}

	_, err = i.api.Block().Put(ctx, bytes.NewReader(block.RawData()),
		options.Block.Hash(mhd.Code, mhd.Length),
		options.Block.Format(multihash.Codes[block.Cid().Type()]))
	return err
}

func (i *IPFSstore) PutMany(ctx context.Context, blocks []blocks.Block) error {
	for _, block := range blocks {
		if err := i.Put(ctx, block); err != nil {
			return err
		}
	}

	return nil
}

func (i *IPFSstore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	return nil, errors.New("not supported")
}

func (i *IPFSstore) DeleteBlock(ctx context.Context, cid cid.Cid) error {
	return errors.New("not supported")
}

func (i *IPFSstore) HashOnRead(enabled bool) {
	return
}
