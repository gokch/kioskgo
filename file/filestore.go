package file

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
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

	if writer != nil {
		err = files.WriteTo(writer.Node, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileStore) PutCid(path string, ci cid.Cid) error {
	fileName := filepath.Join(f.rootPath, path)
	filePath := filepath.Dir(fileName)
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		return err
	}

	// write cids
	err = os.WriteFile(filepath.Join(filePath, DEF_PATH_CID_INFO), ci.Bytes(), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileStore) Get(path string) (*Reader, error) {
	fullPath := filepath.Join(f.rootPath, path)
	reader := NewReaderFromPath(fullPath)
	if reader == nil {
		return nil, os.ErrNotExist
	}
	return reader, nil
}

func (f *FileStore) GetCid(path string) (cid.Cid, error) {
	fullPath := filepath.Join(f.rootPath, path)

	var ci cid.Cid
	rawCid, err := os.ReadFile(filepath.Join(filepath.Dir(fullPath), DEF_PATH_CID_INFO))
	if err == nil {
		_, ci, _ = cid.CidFromBytes(rawCid)
	}
	return ci, nil
}

func (f *FileStore) Exist(path string) bool {
	fullPath := filepath.Join(f.rootPath, path)
	_, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return true
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

func (f *FileStore) Delete(path string) error {
	fullPath := filepath.Join(f.rootPath, path)
	return os.Remove(fullPath)
}

// cids
const (
	DEF_PATH_CID_INFO = "/.info"
)
