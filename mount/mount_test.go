package mount

/*
import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestP2P(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start uploder
	Uploader, err := NewP2P(ctx, "", "oripath", nil)
	require.NoError(t, err)

	cid, err := Uploader.Upload(ctx, "kokomi.png")
	require.NoError(t, err)

	fullAddr := getHostAddress(Uploader.host)
	fmt.Println("addr, cid : ", fullAddr, "|", cid.String())

	// start downloader
	Downloader, err := NewP2P(ctx, "", "cpypath", nil)
	require.NoError(t, err)
	err = Downloader.Connect(ctx, fullAddr)
	require.NoError(t, err)

	// download file
	err = Downloader.Download(ctx, cid, "kokomi.png")
	require.NoError(t, err)

	Uploader.Close()
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
