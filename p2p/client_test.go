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

	fs1 := file.NewFileStore("rootpath1")
	server, err := NewP2PServer(ctx, "", fs1)
	require.NoError(t, err)

	cid, err := server.Upload(ctx, file.NewReaderFromBytes([]byte("test")))
	defer server.Close()

	fullAddr := GetHostAddress(server.host)

	fs2 := file.NewFileStore("rootpath2")
	client, err := NewP2PClient(ctx, "", fs2)
	require.NoError(t, err)

	err = client.Connect(ctx, fullAddr)
	require.NoError(t, err)

	err = client.Download(ctx, cid, "test/abc/d/e.txt")
	require.NoError(t, err)

	client.Disconnect()

}
