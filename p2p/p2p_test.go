package p2p

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

	cid, err := Uploader.Upload(ctx, "맹구.png")
	fullAddr := getHostAddress(Uploader.host)
	fmt.Println(fullAddr, cid.String())

	// start downloader
	Downloader, err := NewP2P(ctx, "", "cpypath", nil)
	require.NoError(t, err)
	err = Downloader.Connect(ctx, fullAddr)
	require.NoError(t, err)

	// download file
	err = Downloader.Download(ctx, cid, "맹구.png")
	require.NoError(t, err)

	Uploader.Close()
}
