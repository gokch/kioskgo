package file

import (
	"context"
	"time"

	theine "github.com/Yiling-J/theine-go"
	datastore "github.com/ipfs/go-datastore"
	query "github.com/ipfs/go-datastore/query"
)

// Cache is a cache that provides methods for storing and retrieving data.
type Cache struct {
	// cache is a theine cache that stores the data.
	cache *theine.Cache[string, []byte]
	// ttl is the time to live for the data in the cache.
	ttl time.Duration
}

// NewCacheStore creates a new cache with the given ttl.
func NewCacheStore(ttl time.Duration) *Cache {
	cache, _ := theine.NewBuilder[string, []byte](1024 * 1024 * 1024).Build()
	return &Cache{
		cache: cache,
		ttl:   ttl,
	}
}

// Put stores the given key and value in the cache.
func (ds *Cache) Put(ctx context.Context, key datastore.Key, value []byte) error {
	ds.cache.SetWithTTL(key.String(), value, 0, ds.ttl)
	return nil
}

// Sync synchronizes the cache with the underlying datastore.
func (ds *Cache) Sync(ctx context.Context, prefix datastore.Key) error {
	return nil
}

// Get retrieves the value for the given key from the cache.
func (ds *Cache) Get(ctx context.Context, key datastore.Key) (value []byte, err error) {
	value, success := ds.cache.Get(key.String())
	if !success {
		return nil, datastore.ErrNotFound
	}
	return value, nil
}

// Has checks if the cache contains the given key.
func (ds *Cache) Has(ctx context.Context, key datastore.Key) (exists bool, err error) {
	val, _ := ds.cache.Get(key.String())
	return val != nil, nil
}

// GetSize retrieves the size of the value for the given key from the cache.
func (ds *Cache) GetSize(ctx context.Context, key datastore.Key) (size int, err error) {
	value, _ := ds.cache.Get(key.String())
	if value == nil {
		return -1, datastore.ErrNotFound
	}
	return len(value), nil
}

// Delete deletes the value for the given key from the cache.
func (ds *Cache) Delete(ctx context.Context, key datastore.Key) (err error) {
	ds.cache.Delete(key.String())
	return nil
}

// Query executes the given query on the cache.
func (ds *Cache) Query(ctx context.Context, q query.Query) (query.Results, error) {
	var keys = make([]string, 0, 1024)
	var vals = make([][]byte, 0, 1024)
	ds.cache.Range(func(key string, value []byte) bool {
		keys = append(keys, key)
		vals = append(vals, value)
		return true
	})
	entries := query.ResultEntriesFrom(keys, vals)
	return query.ResultsWithEntries(q, entries), nil
}

// Batch returns a datastore.Batch that can be used to batch operations on the cache.
func (ds *Cache) Batch(ctx context.Context) (datastore.Batch, error) {
	return nil, datastore.ErrBatchUnsupported
}

// Close closes the cache.
func (ds *Cache) Close() error {
	ds.cache.Close()
	return nil
}
