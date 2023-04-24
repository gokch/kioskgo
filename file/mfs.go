package file

import (
	"context"

	bserv "github.com/ipfs/boxo/blockservice"
	bstore "github.com/ipfs/boxo/blockstore"
	dag "github.com/ipfs/boxo/ipld/merkledag"
	ft "github.com/ipfs/boxo/ipld/unixfs"
	"github.com/ipfs/boxo/mfs"
	dssync "github.com/ipfs/go-datastore/sync"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
)

// TODO : DAG 부터 좀 해놓고 생각하자..
type Mfs struct {
	*mfs.Root
}

func NewMfs(fs *FileStore) *Mfs {
	db := dssync.MutexWrap(fs)
	bs := bstore.NewBlockstore(db)
	blockserv := bserv.New(bs, offline.Exchange(bs))

	root, err := mfs.NewRoot(context.Background(), dag.NewDAGService(blockserv), dag.NodeWithData(ft.FolderPBData()), nil)
	if err != nil {
		return nil
	}
	return &Mfs{
		Root: root,
	}
}
