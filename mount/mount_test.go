package mount

import (
	"bytes"
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

	// make block store
	mds := datastore.NewMapDatastore()
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(mds))
	ex := offline.Exchange(bs)

	// start uploder
	mnt := NewMount(fs, bs)

	dag, err := NewDag(ctx, mnt, ex)
	require.NoError(t, err)

	ci, err := dag.Upload(ctx, "a/a.txt", bytes.NewReader(bytes.Repeat([]byte("abcdfaefedefede"), 100000)))
	require.NoError(t, err)

	fmt.Println(ci.String())

	// // get data from cid
	// cida, err := fs.GetCid("a/a.txt")
	// require.NoError(t, err)

	// newa, err := mnt.Dag.Dagserv.Get(ctx, cida)
	// require.NoError(t, err)

	// cidb, err := fs.GetCid("b/b.txt")
	// require.NoError(t, err)

	// newb, err := mnt.Dag.Dagserv.Get(ctx, cidb)
	// require.NoError(t, err)

	// require.Equal(t, oria, newa.RawData())
	// require.Equal(t, orib, newb.RawData())

	// scida := cida.String()
	// newcida, _ := cid.Decode(scida)
	// require.Equal(t, cida, newcida)
}

/*
func TestP2PCar(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start uploder
	Uploader, err := NewP2P(ctx, "", "oripath", nil)
	require.NoError(t, err)

	cid, err := Uploader.Upload(ctx, "kokomi.png")
	require.NoError(t, err)

	fullAddr := getHostAddress(Uploader.host)
	fmt.Println("addr, cid : ", fullAddr, "|", cid.String())

	err = Uploader.SaveCar(ctx)
	require.NoError(t, err)

	Downloader, err := NewP2P(ctx, "", "oripath", nil)
	require.NoError(t, err)

	err = Downloader.LoadCar(ctx)
	require.NoError(t, err)

	fmt.Println(reflect.DeepEqual(Downloader.dsrv, Uploader.dsrv))
}
*/
