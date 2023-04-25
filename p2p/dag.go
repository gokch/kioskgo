package p2p

import (
	"context"

	"github.com/gokch/kioskgo/file"

	"github.com/ipfs/boxo/exchange"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	"github.com/multiformats/go-multicodec"

	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	chunker "github.com/ipfs/boxo/chunker"
	"github.com/ipfs/boxo/ipld/merkledag"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/boxo/ipld/unixfs/importer/balanced"
	uih "github.com/ipfs/boxo/ipld/unixfs/importer/helpers"
)

type DagBlock struct {
	dag *uih.DagBuilderParams // use MapDataStore
	fs  *file.FileStore       // FileStore
}

func NewDagBlock(ctx context.Context, rootPath string, rem exchange.Interface) (*DagBlock, error) {
	fs := file.NewFileStore(rootPath)

	// make import params
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(dsync.MutexWrap(ds.NewMapDatastore())))
	dsrv := merkledag.NewDAGService(blockservice.New(bs, rem))
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

	dagBlock := &DagBlock{
		dag: builder,
		fs:  fs,
	}

	// import blocks in Merkle-DAG from fileStore
	fs.Iterate("", func(path string, reader *file.Reader) {
		dagBlock.Upload(ctx, path)
	})

	return dagBlock, nil
}

func (p *DagBlock) Download(ctx context.Context, ci cid.Cid, path string) error {
	// conn manager 가 살아있을 때만 download
	node, err := p.dag.Dagserv.Get(ctx, ci)
	if err != nil {
		return err
	}

	unixFSNode, err := unixfile.NewUnixfsFile(ctx, p.dag.Dagserv, node)
	if err != nil {
		return err
	}

	return p.fs.Put(path, file.NewWriter(unixFSNode, node.Cid()))
}

func (p *DagBlock) Upload(ctx context.Context, path string) (cid.Cid, error) {
	data, err := p.fs.Get(path)
	if err != nil {
		return cid.Undef, err
	}
	// Split the file up into fixed sized 256KiB chunks
	ufsBuilder, err := p.dag.New(chunker.NewSizeSplitter(data, chunker.DefaultBlockSize))
	if err != nil {
		return cid.Undef, err
	}
	nd, err := balanced.Layout(ufsBuilder) // Arrange the graph with a balanced layout
	if err != nil {
		return cid.Undef, err
	}

	return nd.Cid(), nil
}
