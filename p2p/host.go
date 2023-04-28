package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
)

/*
	func makeHost(address string, listenPort int) (host host.Host, err error) {
		var opts []libp2p.Option
		if address != "" { // connect to existing host
			peerAddr, err := peer.AddrInfoFromString(address)
			if err != nil {
				return nil, err
			}
			opts = []libp2p.Option{
				libp2p.ListenAddrs(peerAddr.Addrs...),
			}
		} else { // generate new host
			priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
			if err != nil {
				return nil, err
			}

			opts = []libp2p.Option{
				libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)), // port we are listening on, limiting to a single interface and protocol for simplicity
				libp2p.Identity(priv),
			}
		}

		host, err = libp2p.New(opts...)
		if err != nil {
			return nil, err
		}

		return host, nil
	}
*/
func getHostAddress(h host.Host) string {
	addrInfo := host.InfoFromHost(h)
	addr := addrInfo.Addrs[0]
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", addrInfo.ID.String()))

	return addr.Encapsulate(hostAddr).String()
}
