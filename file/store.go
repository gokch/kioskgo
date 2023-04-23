package file

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ipfs/boxo/files"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

type FileStore struct {
	mtx      sync.Mutex
	rootPath string
}

var _ ds.Datastore = (*FileStore)(nil)
var _ ds.Batching = (*FileStore)(nil)
var _ ds.PersistentDatastore = (*FileStore)(nil)

func NewFileStore(rootPath string) *FileStore {
	os.MkdirAll(rootPath, 0755)

	return &FileStore{
		rootPath: rootPath,
		mtx:      sync.Mutex{},
	}
}

func (d *FileStore) KeyFilename(path ds.Key) string {
	return filepath.Join(d.rootPath, path.String(), DEF_PATH_KEY)
}

func (f *FileStore) Overwrite(ctx context.Context, path ds.Key, value []byte) error {
	if exist, _ := f.Has(ctx, path); exist {
		err := f.Delete(ctx, path)
		if err != nil {
			return err
		}
	}

	return f.Put(ctx, path, value)
}

// Put stores the given value.
func (d *FileStore) Put(ctx context.Context, path ds.Key, value []byte) (err error) {
	fileName := d.KeyFilename(path)

	// mkdirall above.
	err = os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, value, 0666)
}

// Sync would ensure that any previous Puts under the prefix are written to disk.
// However, they already are.
func (d *FileStore) Sync(ctx context.Context, prefix ds.Key) error {
	return nil
}

func (f *FileStore) Get(ctx context.Context, path ds.Key) ([]byte, error) {
	fn := f.KeyFilename(path)
	if !isFile(fn) {
		return nil, ds.ErrNotFound
	}

	return ioutil.ReadFile(fn)
}

// Has returns whether the datastore has a value for a given key
func (f *FileStore) Has(ctx context.Context, key ds.Key) (exists bool, err error) {
	return ds.GetBackedHas(ctx, f, key)
}

func (f *FileStore) GetSize(ctx context.Context, key ds.Key) (size int, err error) {
	return ds.GetBackedSize(ctx, f, key)
}

func (f *FileStore) Iterate(path string, fn func(fpath string, value []byte)) error {
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
			bt, _ := ioutil.ReadAll(rf)
			fn(fpath, bt)
		}
		return nil
	})
}

func (f *FileStore) Delete(ctx context.Context, path ds.Key) error {
	fullPath := f.KeyFilename(path)
	if !isFile(fullPath) {
		return nil
	}
	err := os.Remove(fullPath)
	if os.IsNotExist(err) {
		err = nil // idempotent
	}
	return err
}

// Query implements Datastore.Query
func (f *FileStore) Query(ctx context.Context, q query.Query) (query.Results, error) {
	results := make(chan query.Result)

	walkFn := func(path string, info os.FileInfo, _ error) error {
		// remove ds path prefix
		relPath, err := filepath.Rel(f.rootPath, path)
		if err == nil {
			path = filepath.ToSlash(relPath)
		}

		if !info.IsDir() {
			path = strings.TrimSuffix(path, DEF_PATH_KEY)
			var result query.Result
			key := ds.NewKey(path)
			result.Entry.Key = key.String()
			if !q.KeysOnly {
				result.Entry.Value, result.Error = f.Get(ctx, key)
			}
			results <- result
		}
		return nil
	}

	go func() {
		filepath.Walk(f.rootPath, walkFn)
		close(results)
	}()
	r := query.ResultsWithChan(q, results)
	r = query.NaiveQueryApply(q, r)
	return r, nil
}

func (f *FileStore) Close() error {
	return nil
}

func (f *FileStore) Batch(ctx context.Context) (ds.Batch, error) {
	return ds.NewBasicBatch(f), nil
}

// DiskUsage returns the disk size used by the datastore in bytes.
func (f *FileStore) DiskUsage(ctx context.Context) (uint64, error) {
	var du uint64
	err := filepath.Walk(f.rootPath, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if f != nil && f.Mode().IsRegular() {
			du += uint64(f.Size())
		}
		return nil
	})
	return du, err
}
