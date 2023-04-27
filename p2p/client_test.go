package p2p

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := NewClient(ctx, "", "./test")
	require.NoError(t, err)

	ci, err := client.mount.Upload(ctx, "./test/test.txt")
	require.NoError(t, err)

	fmt.Println("connect | address | cid :", client.Self(), ci.String())

	time.Sleep(time.Second * 100)
}
