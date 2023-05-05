package mount

import (
	"context"
	"fmt"
	"testing"

	"github.com/gokch/kioskgo/file"
	"github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/stretchr/testify/require"
)

func TestInitDHT(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// make file store
	fs := file.NewFileStore("rootpath")
	fm := file.NewFileManager()

	// make block store
	mds := datastore.NewMapDatastore()
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(mds))
	ex := offline.Exchange(bs)

	// start uploder
	mnt := NewMount(fs, fm, bs)

	dag, err := NewDag(ctx, -1, mnt, ex)
	require.NoError(t, err)

	// ci, err := dag.Upload(ctx, "a/a.txt", bytes.NewReader(bytes.Repeat([]byte("abcdfaefedefede"), 1000000)))
	ci, err := dag.Upload(ctx, "a/kokomi.png", nil)
	require.NoError(t, err)

	fmt.Println(ci.String())

	err = dag.Download(ctx, ci, "b/kokomi.png")
	require.NoError(t, err)

}
