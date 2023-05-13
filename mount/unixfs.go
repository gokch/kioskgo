package mount

import (
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/boxo/ipld/unixfs"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multicodec"
)

var (
	cidBuilder = cid.V1Builder{ // Use CIDv1 for all links
		Codec:    uint64(multicodec.DagPb),
		MhType:   uint64(multicodec.Sha3_256), // Use SHA3-256 as the hash function
		MhLength: -1,                          // Use the default hash length for the given hash function (in this case 256 bits)
	}
)

func UnixFsFromFiles(node files.Node) (*merkledag.ProtoNode, error) {

	dir := unixfs.EmptyDirNode()
	err := dir.SetCidBuilder(cidBuilder)
	if err != nil {
		return nil, err
	}
	// ??
	return merkledag.NodeWithData(nil), nil
}

func UnixDirFromFiles(node files.Node, dirNode *merkledag.ProtoNode) error {
	switch n := node.(type) {
	case *files.Symlink:
		// TODO.. 할게 없을 것  같은데?
		return nil
	case files.File:
		_ = n
		merkledag.NodeWithData(nil)
		dirNode.AddNodeLink("", nil)
	case files.Directory:

	}

	// dir.AddNodeLink()

	return nil
}
