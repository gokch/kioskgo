package mount

import (
	"context"

	"github.com/gokch/ipfs_mount/file"
	blockstore "github.com/ipfs/boxo/blockstore"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	dsq "github.com/ipfs/go-datastore/query"
	posinfo "github.com/ipfs/go-ipfs-posinfo"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log"
)

var logger = logging.Logger("mount")

type Mount struct {
	fs *file.FileStore
	fm *file.FileManager
	bs blockstore.Blockstore
}

var _ blockstore.Blockstore = (*Mount)(nil)

func NewMount(fs *file.FileStore, fm *file.FileManager, bs blockstore.Blockstore) *Mount {
	return &Mount{
		fs: fs,
		fm: fm,
		bs: bs,
	}
}

// AllKeysChan returns a channel from which to read the keys stored in
// the blockstore. If the given context is cancelled the channel will be closed.
func (f *Mount) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	ctx, cancel := context.WithCancel(ctx)

	a, err := f.bs.AllKeysChan(ctx)
	if err != nil {
		cancel()
		return nil, err
	}

	out := make(chan cid.Cid, dsq.KeysOnlyBufSize)
	go func() {
		defer cancel()
		defer close(out)

		var done bool
		for !done {
			select {
			case c, ok := <-a:
				if !ok {
					done = true
					continue
				}
				select {
				case out <- c:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}

		// Can't do these at the same time because the abstractions around
		// leveldb make us query leveldb for both operations. We apparently
		// cant query leveldb concurrently
		b, err := f.fm.AllKeysChan(ctx)
		if err != nil {
			logger.Error("error querying filestore: ", err)
			return
		}

		done = false
		for !done {
			select {
			case c, ok := <-b:
				if !ok {
					done = true
					continue
				}
				select {
				case out <- c:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

func (f *Mount) DeleteBlock(ctx context.Context, c cid.Cid) error {
	err1 := f.bs.DeleteBlock(ctx, c)
	if err1 != nil && !ipld.IsNotFound(err1) {
		return err1
	}

	info := f.fm.Get(c)
	if info == nil {
		return nil
	}
	err2 := f.fs.Delete(ctx, info.Path)

	// if we successfully removed something from the blockstore, but the
	// filestore didnt have it, return success
	if err2 != nil {
		return err2
	}

	if ipld.IsNotFound(err1) {
		return err1
	}

	return nil
}

func (f *Mount) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	blk, err := f.bs.Get(ctx, c)
	if ipld.IsNotFound(err) {
		info := f.fm.Get(c)
		if info == nil {
			return nil, ipld.ErrNotFound{}
		}

		reader, err := f.fs.Get(ctx, info.Path)
		if err != nil {
			return nil, err
		}
		rawblk, err := reader.GetBlock(info.Offset, info.Size)
		if err != nil {
			return nil, err
		}
		reader.Close()

		blk, err = blocks.NewBlockWithCid(rawblk, c)
		if err != nil {
			return nil, err
		}
		// set block from fs
		err = f.bs.Put(ctx, blk)
		if err != nil {
			return nil, err
		}
		return blk, nil
	}
	return blk, err
}

func (f *Mount) GetSize(ctx context.Context, c cid.Cid) (int, error) {
	size, err := f.bs.GetSize(ctx, c)
	if err != nil {
		if ipld.IsNotFound(err) {
			fi := f.fm.Get(c)
			return fi.Size, nil
		}
		return -1, err
	}
	return size, nil
}

func (f *Mount) Has(ctx context.Context, c cid.Cid) (bool, error) {
	has, err := f.bs.Has(ctx, c)
	if err != nil {
		return false, err
	}

	if has {
		return true, nil
	}

	return f.fm.Has(ctx, c), nil
}

func (f *Mount) Put(ctx context.Context, b blocks.Block) error {
	has, err := f.Has(ctx, b.Cid())
	if err != nil {
		return err
	} else if has {
		return nil
	}

	switch b := b.(type) {
	case *posinfo.FilestoreNode:
		// do nothing. fs will be updated in dag after all block put
	default:
		return f.bs.Put(ctx, b)
	}
	return nil
}

// PutMany is like Put(), but takes a slice of blocks, allowing
// the underlying blockstore to perform batch transactions.
func (f *Mount) PutMany(ctx context.Context, bs []blocks.Block) error {
	var normals []blocks.Block
	for _, b := range bs {
		has, err := f.Has(ctx, b.Cid())
		if err != nil {
			return err
		} else if has {
			continue
		}

		switch b := b.(type) {
		case *posinfo.FilestoreNode:
			// do nothing. fs will be updated in dag after all block put
		default:
			normals = append(normals, b)
		}
	}

	if len(normals) > 0 {
		err := f.bs.PutMany(ctx, normals)
		if err != nil {
			return err
		}
	}
	return nil
}

// HashOnRead calls blockstore.HashOnRead.
func (f *Mount) HashOnRead(enabled bool) {
	f.bs.HashOnRead(enabled)
}
