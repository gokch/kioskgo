package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"
	bsclient "github.com/ipfs/boxo/bitswap/client"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	"github.com/ipfs/boxo/files"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	routinghelpers "github.com/libp2p/go-libp2p-routing-helpers"
	"github.com/multiformats/go-multiaddr"
)

type P2PClient struct {
	Address  string
	host     host.Host
	bsClient *bsclient.Client

	fileStore file.FileStore
}

func NewP2PClient(ctx context.Context, address string, fileStore file.FileStore) (*P2PClient, error) {
	host, address, err := makeHost(address, 3000)
	if err != nil {
		return nil, err
	}
	bitSwapNetwork := bsnet.NewFromIpfsHost(host, routinghelpers.Null{})
	bswap := bsclient.New(ctx, bitSwapNetwork, blockstore.NewBlockstore(datastore.NewNullDatastore()))
	bitSwapNetwork.Start(bswap)

	return &P2PClient{
		Address:   address,
		host:      host,
		bsClient:  bswap,
		fileStore: fileStore,
	}, nil
}

func (p *P2PClient) Close() error {
	if err := p.bsClient.Close(); err != nil {
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

func (p *P2PClient) Disconnect(ctx context.Context, targetPeer string) error {
	// connmanager 가 살아있을 때만 disconnect
	return p.host.ConnManager().Close()
}

func (p *P2PClient) Download(ctx context.Context, ci cid.Cid, path string) error {
	// conn manager 가 살아있을 때만 download
	dserv := merkledag.NewReadOnlyDagService(
		merkledag.NewSession(ctx, merkledag.NewDAGService(blockservice.New(blockstore.NewBlockstore(datastore.NewNullDatastore()), p.bsClient))))
	node, err := dserv.Get(ctx, ci)
	if err != nil {
		return err
	}

	unixFSNode, err := unixfile.NewUnixfsFile(ctx, dserv, node)
	if err != nil {
		return err
	}

	return p.fileStore.Put(path, unixFSNode)
}

func (p *P2PClient) Upload(ctx context.Context, ci cid.Cid, path, name string, data *files.ReaderFile) error {
	return nil
}
