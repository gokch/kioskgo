package file

import (
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

// read from specific path using boxo/files
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
	fileName := filepath.Join(f.rootPath, path)
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (f *FileStore) Delete(path string) error {
	fileName := filepath.Join(f.rootPath, path)
	return os.Remove(fileName)
}
