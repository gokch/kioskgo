package mount

import (
	"fmt"
	"testing"

	"github.com/ipfs/boxo/ipld/unixfs"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multicodec"
	"github.com/stretchr/testify/require"
)

func TestDir(t *testing.T) {
	dirNode := unixfs.EmptyDirNode()
	err := dirNode.SetCidBuilder(cid.V1Builder{ // Use CIDv1 for all links
		Codec:    uint64(multicodec.DagPb),
		MhType:   uint64(multicodec.Sha3_256), // Use SHA3-256 as the hash function
		MhLength: -1,                          // Use the default hash length for the given hash function (in this case 256 bits)
	})
	require.NoError(t, err)

	// fs := file.NewFileStore("rootpath")
	// reader, err := fs.Get(context.Background(), "picture/1.jpg")
	require.NoError(t, err)

	// merkledag.NodeWithData()
	// dirNode.AddNodeLink("1", reader.ReaderFile)

	fmt.Println(dirNode.Cid().String())

}
