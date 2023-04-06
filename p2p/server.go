package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	routinghelpers "github.com/libp2p/go-libp2p-routing-helpers"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multicodec"

	bsnet "github.com/ipfs/boxo/bitswap/network"
	bsserver "github.com/ipfs/boxo/bitswap/server"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	chunker "github.com/ipfs/boxo/chunker"
	offline "github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/boxo/ipld/unixfs/importer/balanced"
	uih "github.com/ipfs/boxo/ipld/unixfs/importer/helpers"
)

type P2PServer struct {
	Address string
	host    host.Host
	bss     *bsserver.Server
	bsn     bsnet.BitSwapNetwork
	builder uih.DagBuilderParams

	fs *file.FileStore
}

func NewP2PServer(ctx context.Context, address string, fs *file.FileStore) (*P2PServer, error) {
	// make import params
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(dsync.MutexWrap(datastore.NewMapDatastore()))) // handle identity multihashes, these don't require doing any actual lookups
	dsrv := merkledag.NewDAGService(blockservice.New(bs, offline.Exchange(bs)))
	// Create a UnixFS graph from our file, parameters described here but can be visualized at https://dag.ipfs.tech/
	params := uih.DagBuilderParams{
		Maxlinks:  uih.DefaultLinksPerBlock, // Default max of 174 links per block
		RawLeaves: true,                     // Leave the actual file bytes untouched instead of wrapping them in a dag-pb protobuf wrapper
		CidBuilder: cid.V1Builder{ // Use CIDv1 for all links
			Codec:    uint64(multicodec.DagPb),
			MhType:   uint64(multicodec.Sha2_256), // Use SHA2-256 as the hash function
			MhLength: -1,                          // Use the default hash length for the given hash function (in this case 256 bits)
		},
		Dagserv: dsrv,
		NoCopy:  false,
	}

	// Start listening on the Bitswap protocol
	// For this example we're not leveraging any content routing (DHT, IPNI, delegated routing requests, etc.) as we know the peer we are fetching from
	host, err := makeHost(address, 0)
	if err != nil {
		return nil, err
	}
	bsn := bsnet.NewFromIpfsHost(host, routinghelpers.Null{})
	bss := bsserver.New(ctx, bsn, bs)

	return &P2PServer{
		Address: getHostAddress(host),
		host:    host,
		bsn:     bsn,
		bss:     bss,
		builder: params,
		fs:      fs,
	}, nil
}

func (p *P2PServer) Start() error {
	p.bsn.Start(p.bss)
	return nil
}

func (p *P2PServer) Close() error {
	p.bsn.Stop()
	return nil
}

func (p *P2PServer) Upload(ctx context.Context, reader *file.Reader) (cid.Cid, error) {
	ufsBuilder, err := p.builder.New(chunker.NewSizeSplitter(reader.ReaderFile, chunker.DefaultBlockSize)) // Split the file up into fixed sized 256KiB chunks
	if err != nil {
		return cid.Undef, err
	}
	nd, err := balanced.Layout(ufsBuilder) // Arrange the graph with a balanced layout
	if err != nil {
		return cid.Undef, err
	}
	return nd.Cid(), nil
}
