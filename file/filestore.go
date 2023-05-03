package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/boxo/files"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	posinfo "github.com/ipfs/go-ipfs-posinfo"
)

type FileStore struct {
	rootPath string
	FM       *FileManager
}

func NewFileStore(rootPath string, blockSize uint64) *FileStore {
	os.MkdirAll(rootPath, 0755)

	return &FileStore{
		rootPath: rootPath,
		FM:       NewFileManager(blockSize),
	}
}

/*
	func (f *FileStore) Overwrite(path string, writer *Writer) error {
		if f.Exist(path) {
			err := f.Delete(path)
			if err != nil {
				return err
			}
		}

		return f.Put(path, writer)
	}
*/

func (f *FileStore) Put(ctx context.Context, ci cid.Cid, posInfo posinfo.PosInfo, writer *Writer) error {
	fileInfo := f.FM.Get(ci)
	f.FM.Put(fileInfo.ci, fileInfo.path, fileInfo.offset) // put offset

	fileName := filepath.Join(f.rootPath, fileInfo.path)
	filePath := filepath.Dir(fileName)
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		return err
	}

	if writer != nil {
		err = files.WriteTo(writer.Node, fileName)
		if err != nil {
			return err
		}
	}

	// write cids
	err = os.WriteFile(filepath.Join(filePath, DEF_PATH_CID_INFO), ci.Bytes(), 0755)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStore) Get(ctx context.Context, ci cid.Cid) (blocks.Block, error) {
	info := f.FM.Get(ci)
	if info == nil {
		return nil, os.ErrNotExist
	}

	fullPath := filepath.Join(f.rootPath, info.path)
	reader := NewReaderFromPath(fullPath)
	if reader == nil {
		return nil, os.ErrNotExist
	}

	return reader.GetBlock(int64(info.offset), int64(info.size), ci)
}

func (f *FileStore) GetSize(ctx context.Context, c cid.Cid) (int, error) {
	info := f.FM.Get(c)
	if info == nil {
		return -1, os.ErrNotExist
	}
	return int(info.size), nil
}

func (f *FileStore) Has(ctx context.Context, ci cid.Cid) (bool, error) {
	info := f.FM.Get(ci)
	if info != nil {
		return true, nil
	}
	return false, nil
}

func (f *FileStore) Iterate(path string, fn func(fpath string, reader *Reader) error) error {
	fullPath := filepath.Join(f.rootPath, path)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	sf, err := files.NewSerialFile(fullPath, true, stat)
	if err != nil {
		return err
	}
	return files.Walk(sf, func(fpath string, node files.Node) error {
		if rf, ok := node.(*files.ReaderFile); ok {
			defer rf.Close()
			if err = fn(fpath, NewReader(rf)); err != nil {
				return err
			}
		}
		return nil // ignore directory
	})
}

func (f *FileStore) DeleteBlock(ctx context.Context, ci cid.Cid) error {
	info := f.FM.Get(ci)
	if info == nil {
		return os.ErrNotExist
	}

	fullPath := filepath.Join(f.rootPath, info.path)
	return os.Remove(fullPath)
}

func (f *FileStore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	out := make(chan cid.Cid)
	for ci := range f.FM.cids {
		out <- ci
	}
	return out, nil
}

// cids
const (
	DEF_PATH_CID_INFO = "/.info"
)
