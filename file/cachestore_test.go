package file

import (
	"context"
	"testing"
	"time"

	"github.com/ipfs/go-datastore"
	"github.com/stretchr/testify/require"
)

func TestCacheStore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := NewCacheStore(time.Second * 300)

	err := cs.Put(ctx, datastore.NewKey("test/abc/d/e.txt"), ([]byte("test")))
	require.NoError(t, err)

	value, err := cs.Get(ctx, datastore.NewKey("test/abc/d/e.txt"))
	require.NoError(t, err)

	require.Equal(t, []byte("test"), value)
}
