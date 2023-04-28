package p2p

import (
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/multiformats/go-multiaddr"
)

func makeHost(address string, listenPort int) (host host.Host, err error) {
	var opts []libp2p.Option

	cm, err := connmgr.NewConnManager(1, 100, connmgr.WithGracePeriod(0))
	if err != nil {
		return nil, err
	}

	if address != "" { // connect to existing host
		peerAddr, err := peer.AddrInfoFromString(address)
		if err != nil {
			return nil, err
		}

		opts = []libp2p.Option{
			libp2p.ConnectionManager(cm),
			libp2p.ListenAddrs(peerAddr.Addrs...),
		}
	} else { // generate new host
		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, rand.Reader)
		if err != nil {
			return nil, err
		}

		opts = []libp2p.Option{
			libp2p.ConnectionManager(cm),
			libp2p.Identity(priv),
		}
	}

	host, err = libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func getHostAddress(h host.Host) string {
	addrInfo := host.InfoFromHost(h)
	addr := addrInfo.Addrs[0]
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", addrInfo.ID.String()))

	return addr.Encapsulate(hostAddr).String()
}

func encodeAddrInfo(targetPeer string) (info *peer.AddrInfo, err error) {
	maddr, err := multiaddr.NewMultiaddr(targetPeer)
	if err != nil {
		return
	}

	// Extract the peer ID from the multiaddr.
	info, err = peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return
	}
	return info, nil
}
