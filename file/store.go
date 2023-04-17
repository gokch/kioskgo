package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/ipfs/boxo/files"
)

type FileStore struct {
	mtx      *sync.Mutex
	rootPath string
}

func NewFileStore(rootPath string) *FileStore {
	os.MkdirAll(rootPath, 0755)

	return &FileStore{
		rootPath: rootPath,
		mtx:      &sync.Mutex{},
	}
}

func (f *FileStore) Overwrite(path string, writer *Writer) error {
	exist, err := f.Exist(path)
	if err != nil {
		return err
	} else if exist {
		err = f.Delete(path)
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
	fileName := filepath.Join(f.rootPath, path)
	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}

	node, err := files.NewSerialFile(fileName, true, stat)
	if err != nil {
		return nil, err
	}
	return NewReader(node.(*files.ReaderFile)), nil
}

func (f *FileStore) Exist(path string) (bool, error) {
	_, _, err := f.makePath(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (f *FileStore) Iterate(path string) ([]*Reader, error) {
	fullPath, stat, err := f.makePath(path)
	if err != nil {
		return nil, err
	}

	sf, err := files.NewSerialFile(fullPath, true, stat)
	if err != nil {
		return nil, err
	}
	readers := make([]*Reader, 0, 100)
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
	fullPath, _, err := f.makePath(path)
	if err != nil {
		return err
	}
	return os.Remove(fullPath)
}

func (f *FileStore) makePath(paths ...string) (string, fs.FileInfo, error) {
	// append root path
	paths = append([]string{f.rootPath}, paths...)
	fullPath := filepath.Join(paths...)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return "", nil, err
	}
	return fullPath, stat, nil
}
