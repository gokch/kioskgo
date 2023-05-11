package p2p

import (
	"time"

	"github.com/ipfs/boxo/ipns"
	ipns_pb "github.com/ipfs/boxo/ipns/pb"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func AddIPNS(cid *cid.Cid, name string) (*ipns_pb.IpnsEntry, error) {
	privateKey, publicKey, err := crypto.GenerateKeyPair(crypto.Ed25519, 2048)
	if err != nil {
		return nil, err
	}

	// Create an IPNS record that expires in one hour and points to the IPFS address
	// /ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5
	ipnsRecord, err := ipns.Create(privateKey, []byte("/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5"), 0, time.Now().Add(1*time.Hour), time.Minute)
	if err != nil {
		return nil, err
	}
	err = ipns.EmbedPublicKey(publicKey, ipnsRecord)
	if err != nil {
		return nil, err
	}

	return ipnsRecord, nil
}
