package p2p

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
)

func TestHostDht(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	host, err := libp2p.New()
	require.NoError(t, err)

	address := getHostAddress(host)
	fmt.Println(address)

	ipfsdht, err := dht.New(ctx, host, dht.BootstrapPeers(dht.GetDefaultBootstrapPeerAddrInfos()...))
	require.NoError(t, err)

	self := ipfsdht.PeerID().String()
	fmt.Println("self :", self)

	pid, err := peer.Decode("QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN")
	require.NoError(t, err)

	addrInfo, err := ipfsdht.FindPeer(ctx, pid)
	require.NoError(t, err)

	fmt.Println(addrInfo.ID, addrInfo.Addrs)

	_ = ipfsdht
}

func TestClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := NewClient(ctx, "", "./oripath")
	require.NoError(t, err)

	ci, err := client.mount.Upload(ctx, "./kokomi.png")
	require.NoError(t, err)

	fmt.Println("connect | address | cid :", client.Self(), ci.String())

	time.Sleep(time.Second * 10)
}
