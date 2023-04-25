package file

import (
	"testing"

	ds "github.com/ipfs/go-datastore"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestDataStoreNewGet(t *testing.T) {
	ctx := context.Background()

	fs := NewDataStore("rootpath")
	require.NotNil(t, fs)

	data1 := []byte("test1")
	err := fs.Overwrite(ctx, ds.NewKey("test/abc/d/e.txt"), data1)
	require.NoError(t, err)

	data2, err := fs.Get(ctx, ds.NewKey("test/abc/d/e.txt"))
	require.NoError(t, err)

	require.Equal(t, data1, data2)
}

/*
func TestStoreIterate(t *testing.T) {
	ctx := context.Background()

	fs := NewFileStore("rootpath")
	require.NotNil(t, fs)

	// add datas
	data := []byte("test1")
	data2 := []byte("test2")
	data3 := []byte("test3")
	err := fs.Overwrite(ctx, ds.NewKey("a/data1.txt"), data)
	require.NoError(t, err)
	err = fs.Overwrite(ctx, ds.NewKey("a/b/c/data2.txt"), data2)
	require.NoError(t, err)
	err = fs.Overwrite(ctx, ds.NewKey("a/b/c/d/e/data3.txt"), data3)
	require.NoError(t, err)

	// iterate
	err = fs.Iterate("", func(fpath string, value []byte) {
		fmt.Println(fpath, value)
	})
	require.NoError(t, err)
}
*/
