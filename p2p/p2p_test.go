package p2p

import (
	"context"
	"fmt"
	"testing"

	"github.com/gokch/kioskgo/file"
	"github.com/stretchr/testify/require"
)

func TestP2P(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start uploder
	fs1 := file.NewFileStore("")
	Uploader, err := NewP2P(ctx, "", fs1, nil)
	require.NoError(t, err)

	cid, err := Uploader.Upload(ctx, file.NewReaderFromPath("./맹구.png"))
	fullAddr := getHostAddress(Uploader.host)
	fmt.Println(fullAddr, cid.String())

	// start downloader
	fs2 := file.NewFileStore("rootpath")
	Downloader, err := NewP2P(ctx, "", fs2, nil)
	require.NoError(t, err)
	err = Downloader.Connect(ctx, fullAddr)
	require.NoError(t, err)

	// download file
	err = Downloader.Download(ctx, cid, "new/맹구.png")
	require.NoError(t, err)

	Uploader.Close()
}
