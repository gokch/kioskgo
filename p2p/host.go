package p2p

import (
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func makeHost(address string, listenPort int) (host host.Host, addr string, err error) {
	var opts []libp2p.Option
	if address != "" { // connect to existing host
		peerAddr, err := peer.AddrInfoFromString(address)
		if err != nil {
			return nil, "", err
		}
		opts = []libp2p.Option{
			libp2p.ListenAddrs(peerAddr.Addrs...),
		}
	} else { // generate new host
		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, rand.Reader)
		if err != nil {
			return nil, "", err
		}

		opts = []libp2p.Option{
			libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)), // port we are listening on, limiting to a single interface and protocol for simplicity
			libp2p.Identity(priv),
		}
	}

	host, err = libp2p.New(opts...)
	if err != nil {
		return nil, "", err
	}
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", host.ID().String()))
	addr = host.Addrs()[0].Encapsulate(hostAddr).String()

	return host, addr, nil
}
