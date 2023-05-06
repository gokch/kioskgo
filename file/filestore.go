package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/boxo/files"
)

type FileStore struct {
	rootPath string
}

func NewFileStore(rootPath string) *FileStore {
	os.MkdirAll(rootPath, 0755)

	return &FileStore{
		rootPath: rootPath,
	}
}

func (f *FileStore) Overwrite(ctx context.Context, path string, writer *Writer) error {
	err := f.Delete(ctx, path)
	if err != nil {
		return err
	}

	return f.Put(ctx, path, writer)
}

func (f *FileStore) Put(ctx context.Context, path string, writer *Writer) error {
	fileName := filepath.Join(f.rootPath, path)
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

	return nil
}

func (f *FileStore) Get(ctx context.Context, path string) (*Reader, error) {
	fullPath := filepath.Join(f.rootPath, path)
	return NewReaderFromPath(fullPath), nil
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

func (f *FileStore) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(f.rootPath, path)
	return os.Remove(fullPath)
}

// cids
const (
	DEF_PATH_CID_INFO = "/.info"
)
