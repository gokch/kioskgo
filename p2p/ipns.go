package p2p

import (
	"time"

	"github.com/ipfs/boxo/ipns"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func AddIPNS(cid *cid.Cid, name string) (*ipns.Record, error) {
	privateKey, _, err := crypto.GenerateKeyPair(crypto.Ed25519, 2048)
	if err != nil {
		return nil, err
	}

	// Create an IPNS record that expires in one hour and points to the IPFS address
	// /ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5
	p, err := path.NewPath("/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5")
	if err != nil {
		return nil, err
	}
	
	ipnsRecord, err := ipns.NewRecord(privateKey, p, 0, time.Now().Add(1*time.Hour), time.Minute, ipns.WithPublicKey(true))
	if err != nil {
		return nil, err
	}
	// ipnsRecord.

	return ipnsRecord, nil
}
