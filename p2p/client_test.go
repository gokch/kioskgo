package p2p

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-nat"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"

	"github.com/libp2p/go-libp2p/p2p/host/autonat"
)

func TestHostDht(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	host, err := libp2p.New()
	require.NoError(t, err)

	address := getHostAddress(host)
	fmt.Println(address)

	// connect to bootstrap nodes
	bootstraps := dht.GetDefaultBootstrapPeerAddrInfos()
	for _, addrInfo := range bootstraps {
		if err := host.Connect(ctx, addrInfo); err != nil {
			fmt.Printf("failed to connect to bootstrap node %s: %s\n", addrInfo.ID, err)
		}
	}
	fmt.Println(bootstraps)

	ipfsdht, err := dht.New(ctx, host, dht.Mode(dht.ModeServer), dht.BootstrapPeers(dht.GetDefaultBootstrapPeerAddrInfos()...))
	require.NoError(t, err)

	err = ipfsdht.Bootstrap(ctx)
	require.NoError(t, err)

	self := ipfsdht.PeerID().String()
	fmt.Println("self :", self)

	Newpid, err := peer.Decode("QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN")
	require.NoError(t, err)

	err = ipfsdht.Ping(ctx, Newpid)
	require.NoError(t, err)

	ma, err := ipfsdht.FindPeer(ctx, Newpid)
	require.NoError(t, err)

	fmt.Println(ma.String())
}

func TestClient(t *testing.T) {
	// upload
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := NewClient(ctx, "./oripath")
	require.NoError(t, err)

	ci, err := client.mount.Upload(ctx, "kokomi.png")
	require.NoError(t, err)

	fmt.Println("connect | address | cid :", client.Self(), ci.String())

	// download
	client2, err := NewClient(ctx, "./cpypath")
	require.NoError(t, err)

	err = client2.Connect(ctx, client.Self())
	require.NoError(t, err)

	err = client2.mount.Download(ctx, ci, "kokomi.png")
	require.NoError(t, err)

}

func TestIPv6(t *testing.T) {
	maddr, err := multiaddr.NewMultiaddr("/ip6/::/tcp/0")
	require.NoError(t, err)

	fmt.Println(maddr)
}

func TestForwarding(t *testing.T) {
	listenAddr, err := net.ResolveUDPAddr("udp6", "[::]:0")
	require.NoError(t, err)

	listener, err := net.ListenUDP("udp6", listenAddr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Forwarding address: /ip6/%s/udp/%d/quic\n", listener.LocalAddr().(*net.UDPAddr).IP.String(), listener.LocalAddr().(*net.UDPAddr).Port)

}

func TestNat(t *testing.T) {

	ctx := context.Background()
	na, err := nat.DiscoverGateway(ctx)
	require.NoError(t, err)
	fmt.Println("수신됨 :", na)
	if na != nil {
		fmt.Println(na.Type())
	}

	// nat := nat.DiscoverNATs(ctx)

	// select {
	// case n := <-nat:
	// 	fmt.Println("수신됨 :", n)
	// 	if n != nil {
	// 		fmt.Println(n.Type())
	// 	}
	// case <-time.After(time.Second * 3000):
	// 	fmt.Println("timeout")
	// }

}

func TestAutoNat(t *testing.T) {
	host, err := libp2p.New()
	require.NoError(t, err)
	// Autonat 기능을 사용하여 public IP를 가져옵니다.
	autonatService, err := autonat.New(host)
	fmt.Println(autonatService.Status())

	publicAddr, err := autonatService.PublicAddr()
	require.NoError(t, err)

	fmt.Println(publicAddr.MarshalJSON())
}
