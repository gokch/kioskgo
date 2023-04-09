package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	routinghelpers "github.com/libp2p/go-libp2p-routing-helpers"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multicodec"

	"github.com/ipfs/boxo/bitswap"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	chunker "github.com/ipfs/boxo/chunker"
	offline "github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/ipld/merkledag"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/boxo/ipld/unixfs/importer/balanced"
	uih "github.com/ipfs/boxo/ipld/unixfs/importer/helpers"
	"github.com/ipfs/boxo/routing/http/contentrouter"
)

type P2P struct {
	Address string
	host    host.Host
	bsn     bsnet.BitSwapNetwork
	bswap   *bitswap.Bitswap
	builder *uih.DagBuilderParams
	fs      *file.FileStore
}

func NewP2P(ctx context.Context, address string, fs *file.FileStore, clientrouter contentrouter.Client) (*P2P, error) {
	// make import params
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(dsync.MutexWrap(datastore.NewMapDatastore()))) // handle identity multihashes, these don't require doing any actual lookups
	dsrv := merkledag.NewDAGService(blockservice.New(bs, offline.Exchange(bs)))
	// Create a UnixFS graph from our file, parameters described here but can be visualized at https://dag.ipfs.tech/
	builder := &uih.DagBuilderParams{
		Maxlinks:  uih.DefaultLinksPerBlock, // Default max of 174 links per block
		RawLeaves: true,                     // Leave the actual file bytes untouched instead of wrapping them in a dag-pb protobuf wrapper
		CidBuilder: cid.V1Builder{ // Use CIDv1 for all links
			Codec:    uint64(multicodec.DagPb),
			MhType:   uint64(multicodec.Sha3_256), // Use SHA3-256 as the hash function
			MhLength: -1,                          // Use the default hash length for the given hash function (in this case 256 bits)
		},
		Dagserv: dsrv,
		NoCopy:  false,
	}

	host, err := makeHost(address, 0)
	if err != nil {
		return nil, err
	}
	address = getHostAddress(host)

	var r routing.ContentRouting
	if clientrouter == nil {
		r = routinghelpers.Null{}
	} else {
		r = contentrouter.NewContentRoutingClient(clientrouter)
	}
	bsn := bsnet.NewFromIpfsHost(host, r)
	bswap := bitswap.New(ctx, bsn, bs)
	bsn.Start(bswap)

	return &P2P{
		Address: address,
		host:    host,
		bsn:     bsn,
		bswap:   bswap,
		builder: builder,
		fs:      fs,
	}, nil
}

func (p *P2P) Close() error {
	p.bsn.Stop()
	if err := p.bswap.Close(); err != nil {
		return err
	}
	if err := p.host.Close(); err != nil {
		return err
	}
	return nil
}

func (p *P2P) Connect(ctx context.Context, targetPeer string) error {
	// TODO : conn manager 가 죽었을 때만 connect
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

func (p *P2P) Download(ctx context.Context, ci cid.Cid, path string) error {
	// conn manager 가 살아있을 때만 download
	dserv := merkledag.NewReadOnlyDagService(
		merkledag.NewSession(ctx, merkledag.NewDAGService(blockservice.New(blockstore.NewBlockstore(datastore.NewNullDatastore()), p.bswap))))
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

func (p *P2P) Upload(ctx context.Context, reader *file.Reader) (cid.Cid, error) {
	// Split the file up into fixed sized 256KiB chunks
	ufsBuilder, err := p.builder.New(chunker.NewSizeSplitter(reader.ReaderFile, chunker.DefaultBlockSize))
	if err != nil {
		return cid.Undef, err
	}
	nd, err := balanced.Layout(ufsBuilder) // Arrange the graph with a balanced layout
	if err != nil {
		return cid.Undef, err
	}
	return nd.Cid(), nil
}
