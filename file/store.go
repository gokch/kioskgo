package file

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/ipfs/boxo/files"
)

type FileStore struct {
	mtx      sync.Mutex
	rootPath string
}

func NewFileStore(rootPath string) *FileStore {
	os.MkdirAll(rootPath, 0755)

	return &FileStore{
		rootPath: rootPath,
		mtx:      sync.Mutex{},
	}
}

func (f *FileStore) Overwrite(path string, writer *Writer) error {
	if f.Exist(path) {
		err := f.Delete(path)
		if err != nil {
			return err
		}
	}

	return f.Put(path, writer)
}

func (f *FileStore) Put(path string, writer *Writer) error {
	fileName := filepath.Join(f.rootPath, path)
	filePath := filepath.Dir(fileName)
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		return err
	}

	return files.WriteTo(writer.Node, fileName)
}

func (f *FileStore) Get(path string) (*Reader, error) {
	fullPath := filepath.Join(f.rootPath, path)
	return NewReaderFromPath(fullPath), nil
}

func (f *FileStore) Exist(path string) bool {
	fullPath := filepath.Join(f.rootPath, path)
	_, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return true
}

func (f *FileStore) Iterate(path string) ([]*Reader, error) {
	fullPath := filepath.Join(f.rootPath, path)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	sf, err := files.NewSerialFile(fullPath, true, stat)
	if err != nil {
		return nil, err
	}
	readers := make([]*Reader, 0, 1024)
	err = files.Walk(sf, func(fpath string, node files.Node) error {
		if rf, ok := node.(*files.ReaderFile); ok {
			readers = append(readers, NewReader(rf))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return readers, nil
}

func (f *FileStore) Delete(path string) error {
	fullPath := filepath.Join(f.rootPath, path)
	return os.Remove(fullPath)
}
