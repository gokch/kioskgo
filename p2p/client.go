package p2p

import (
	"bytes"
	"context"
	"io"

	bsclient "github.com/ipfs/boxo/bitswap/client"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	files "github.com/ipfs/go-ipfs-files"
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
}

func NewP2PClient(ctx context.Context, address string) (*P2PClient, error) {
	host, address, err := makeHost(address, 3000)
	if err != nil {
		return nil, err
	}
	bitSwapNetwork := bsnet.NewFromIpfsHost(host, routinghelpers.Null{})
	bswap := bsclient.New(ctx, bitSwapNetwork, blockstore.NewBlockstore(datastore.NewNullDatastore()))
	bitSwapNetwork.Start(bswap)

	return &P2PClient{
		Address:  address,
		host:     host,
		bsClient: bswap,
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

func (p *P2PClient) Download(ctx context.Context, ci cid.Cid) ([]byte, error) {
	// conn manager 가 살아있을 때만 download
	dserv := merkledag.NewReadOnlyDagService(
		merkledag.NewSession(ctx, merkledag.NewDAGService(blockservice.New(blockstore.NewBlockstore(datastore.NewNullDatastore()), p.bsClient))))
	node, err := dserv.Get(ctx, ci)
	if err != nil {
		return nil, err
	}

	unixFSNode, err := unixfile.NewUnixfsFile(ctx, dserv, node)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if f, ok := unixFSNode.(files.File); ok {
		if _, err := io.Copy(&buf, f); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
