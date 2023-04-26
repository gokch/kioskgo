package mount

import (
	"context"
	"fmt"
	"testing"

	"github.com/gokch/kioskgo/file"
	"github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/stretchr/testify/require"
)

// file to unixfs ( + cid ) 기능 필요.. file <-> unixfs 필요
func TestInitDHT(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// make file store
	fs := file.NewFileStore("oripath")
	fs.Put("a", file.NewWriterFromBytes([]byte("testaa"), cid.Cid{}))
	fs.Put("b", file.NewWriterFromBytes([]byte("testbb"), cid.Cid{}))

	// make block store
	mds := datastore.NewMapDatastore()
	bs := blockstore.NewIdStore(blockstore.NewBlockstore(mds))
	ex := offline.Exchange(bs)

	// start uploder
	mnt, err := NewMount(ctx, fs, bs, ex)
	require.NoError(t, err)

	// get cid
	protonode := merkledag.NodeWithData([]byte("testaa"))
	protonode.SetCidBuilder(mnt.Dag.CidBuilder)

	cid := protonode.Cid()
	require.NoError(t, err)
	fmt.Println(cid.String())

	fmt.Println(mds)

	// get data from cid
	data, err := mnt.Dag.Dagserv.Get(ctx, cid)
	require.NoError(t, err)
	fmt.Println(data)

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
