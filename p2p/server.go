package p2p

import (
	"context"

	ds_badger "github.com/ipfs/go-ds-badger2"

	"github.com/ipfs/boxo/bitswap"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	"github.com/ipfs/boxo/filestore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

type Server struct {
	// host is the libp2p host.
	host host.Host
	// bswap is the bitswap server.
	bswap *bitswap.Bitswap
}

// NewClient creates a new client.
func NewServer(ctx context.Context, rootPath string) (*Server, error) {
	// init fs
	ds, err := ds_badger.NewDatastore(rootPath, nil)
	if err != nil {
		return nil, err
	}

	bs := blockstore.NewBlockstore(ds)
	fm := filestore.NewFileManager(ds, rootPath)
	fs := filestore.NewFilestore(bs, fm)

	host, err := makeHost(rootPath)
	if err != nil {
		return nil, err
	}

	// init dht for bitswap network
	ipfsdht, err := dht.New(ctx, host)
	if err != nil {
		return nil, err
	}

	// init bitswap network
	bsn := bsnet.NewFromIpfsHost(host, ipfsdht)
	bswap := bitswap.New(ctx, bsn, fs)
	bsn.Start(bswap)

	c := &Server{

		host:  host,
		bswap: bswap,
	}

	return c, nil
}
