package mount

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gokch/kioskgo/file"
	"github.com/ipfs/boxo/exchange/offline"
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
	cs := file.NewCacheStore(time.Second * 300)
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(cs))
	ex := offline.Exchange(bs)

	// start uploder
	mnt := NewMount(fs, fm, bs)

	dag, err := NewDag(ctx, -1, mnt, ex)
	require.NoError(t, err)

	// ci, err := dag.Upload(ctx, "a/a.txt", bytes.NewReader(bytes.Repeat([]byte("abcdfaefedefede"), 1000000)))
	ci, err := dag.Upload(ctx, "picture/1.jpg", nil)
	require.NoError(t, err)

	fmt.Println(ci.String())

	err = dag.Download(ctx, ci, "picture2/1.jpg")
	require.NoError(t, err)

}
