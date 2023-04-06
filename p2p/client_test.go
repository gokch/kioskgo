package p2p

import (
	"context"
	"testing"

	"github.com/gokch/kioskgo/file"
	"github.com/stretchr/testify/require"
)

func TestP2P(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start server
	fs1 := file.NewFileStore("")
	server, err := NewP2PServer(ctx, "", fs1)
	require.NoError(t, err)
	err = server.Start()
	require.NoError(t, err)

	cid, err := server.Upload(ctx, file.NewReaderFromPath("./맹구.png"))
	fullAddr := getHostAddress(server.host)

	// start client
	fs2 := file.NewFileStore("rootpath")
	client, err := NewP2PClient(ctx, "", fs2)
	require.NoError(t, err)
	err = client.Connect(ctx, fullAddr)
	require.NoError(t, err)

	// download file
	err = client.Download(ctx, cid, "new/맹구.png")
	require.NoError(t, err)

	client.Disconnect()
	server.Close()
}
