package file

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

type DataStore struct {
	RootPath string
}

var _ ds.Datastore = (*DataStore)(nil)
var _ ds.Batching = (*DataStore)(nil)
var _ ds.PersistentDatastore = (*DataStore)(nil)

func NewDataStore(rootPath string) *DataStore {
	os.MkdirAll(rootPath, 0755)

	return &DataStore{
		RootPath: rootPath,
	}
}

func (f *DataStore) Overwrite(ctx context.Context, path ds.Key, value []byte) error {
	if exist, _ := f.Has(ctx, path); exist {
		err := f.Delete(ctx, path)
		if err != nil {
			return err
		}
	}

	return f.Put(ctx, path, value)
}

// Put stores the given value.
func (d *DataStore) Put(ctx context.Context, path ds.Key, value []byte) (err error) {
	fileName := getFilename(d.RootPath, path)

	// mkdirall above.
	err = os.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, value, 0666)
}

// Sync would ensure that any previous Puts under the prefix are written to disk.
// However, they already are.
func (d *DataStore) Sync(ctx context.Context, prefix ds.Key) error {
	return nil
}

func (f *DataStore) Get(ctx context.Context, path ds.Key) ([]byte, error) {
	fileName := getFilename(f.RootPath, path)
	if !isFile(fileName) {
		return nil, ds.ErrNotFound
	}

	return os.ReadFile(fileName)
}

// Has returns whether the datastore has a value for a given key
func (f *DataStore) Has(ctx context.Context, path ds.Key) (exists bool, err error) {
	return ds.GetBackedHas(ctx, f, path)
}

func (f *DataStore) GetSize(ctx context.Context, path ds.Key) (size int, err error) {
	return ds.GetBackedSize(ctx, f, path)
}

/*
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
*/
func (f *DataStore) Delete(ctx context.Context, path ds.Key) error {
	fileName := getFilename(f.RootPath, path)

	err := os.Remove(fileName)
	if os.IsNotExist(err) {
		err = nil // idempotent
	}
	return err
}

// Query implements Datastore.Query
func (f *DataStore) Query(ctx context.Context, q query.Query) (query.Results, error) {
	results := make(chan query.Result)

	walkFn := func(path string, info os.FileInfo, _ error) error {
		// remove ds path prefix
		relPath, err := filepath.Rel(f.RootPath, path)
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
		filepath.Walk(f.RootPath, walkFn)
		close(results)
	}()
	r := query.ResultsWithChan(q, results)
	r = query.NaiveQueryApply(q, r)
	return r, nil
}

func (f *DataStore) Close() error {
	return nil
}

func (f *DataStore) Batch(ctx context.Context) (ds.Batch, error) {
	return ds.NewBasicBatch(f), nil
}

// DiskUsage returns the disk size used by the datastore in bytes.
func (f *DataStore) DiskUsage(ctx context.Context) (uint64, error) {
	var du uint64
	err := filepath.Walk(f.RootPath, func(p string, f os.FileInfo, err error) error {
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
