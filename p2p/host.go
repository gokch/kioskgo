package p2p

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/multiformats/go-multiaddr"
)

var (
	// sha256.Sum256([]byte("smpeople"))
	psk = []byte{20, 174, 197, 74, 226, 233, 89, 172, 139, 157, 212, 111, 186, 100, 161, 59, 207, 51, 57, 139, 94, 184, 106, 212, 81, 159, 98, 18, 102, 118, 205, 149}
)

func makeHost(privateKey string) (host host.Host, err error) {
	var priv crypto.PrivKey

	if privateKey == "" {
		priv, _, err = crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, rand.Reader)
	} else {
		privRaw, _ := base64.StdEncoding.DecodeString(privateKey)
		priv, err = crypto.UnmarshalEd25519PrivateKey(privRaw)
	}
	if err != nil {
		return nil, err
	}

	cm, err := connmgr.NewConnManager(1, 1000, connmgr.WithGracePeriod(0))
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ConnectionManager(cm),
		libp2p.PrivateNetwork(psk),
		libp2p.Identity(priv),
		// libp2p.Transport(quic.NewTransport), // QUIC doesn't support private networks yet
	}

	host, err = libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	// raw, _ := priv.Raw()
	// fmt.Println(base64.StdEncoding.EncodeToString(raw))

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
