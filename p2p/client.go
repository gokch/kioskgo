package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"
	bsclient "github.com/ipfs/boxo/bitswap/client"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/ipld/merkledag"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	routinghelpers "github.com/libp2p/go-libp2p-routing-helpers"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type P2PClient struct {
	host host.Host
	bsn  bsnet.BitSwapNetwork
	bsc  *bsclient.Client
	fs   *file.FileStore
}

func NewP2PClient(ctx context.Context, address string, fs *file.FileStore) (*P2PClient, error) {
	host, err := makeHost(address, 0)
	if err != nil {
		return nil, err
	}
	bsn := bsnet.NewFromIpfsHost(host, routinghelpers.Null{})
	bsc := bsclient.New(ctx, bsn, blockstore.NewBlockstore(datastore.NewNullDatastore()))
	bsn.Start(bsc)

	return &P2PClient{
		host: host,
		bsn:  bsn,
		bsc:  bsc,
		fs:   fs,
	}, nil
}

func (p *P2PClient) Close() error {
	p.bsn.Stop()
	if err := p.bsc.Close(); err != nil {
		return err
	}
	if err := p.host.Close(); err != nil {
		return err
	}
	return nil
}

func (p *P2PClient) Connect(ctx context.Context, targetPeer string) error {
	// conn manager 가 죽었을 때만 connect
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
	if err := p.host.Connect(ctx, *info); err != nil {
		return err
	}

	return nil
}

func (p *P2PClient) Disconnect() error {
	// connmanager 가 살아있을 때만 disconnect
	return p.host.ConnManager().Close()
}

func (p *P2PClient) Download(ctx context.Context, ci cid.Cid, path string) error {
	// conn manager 가 살아있을 때만 download
	dserv := merkledag.NewReadOnlyDagService(
		merkledag.NewSession(ctx, merkledag.NewDAGService(blockservice.New(blockstore.NewBlockstore(datastore.NewNullDatastore()), p.bsc))))
	node, err := dserv.Get(ctx, ci)
	if err != nil {
		return err
	}

	unixFSNode, err := unixfile.NewUnixfsFile(ctx, dserv, node)
	if err != nil {
		return err
	}

	return p.fs.Put(path, file.NewWriter(unixFSNode))
}

func (p *P2PClient) Upload(ctx context.Context, ci cid.Cid, path, name string, data *file.Reader) error {

	return nil
}
