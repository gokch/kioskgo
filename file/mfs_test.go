package file

import (
	"context"
	"fmt"
	"testing"

	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/boxo/mfs"
	"github.com/stretchr/testify/require"
)

func TestMFS(t *testing.T) {
	mfsys := NewMfs(NewFileStore("rootpath"))
	rootDir := mfsys.Root.GetDirectory()
	rootDir.Mkdir("test")
	err := mfs.PutNode(mfsys.Root, "test/test.txt", merkledag.NodeWithData([]byte("안녕하세용")))
	require.NoError(t, err)

	child, err := mfs.FlushPath(context.Background(), mfsys.Root, "test")
	require.NoError(t, err)

	// fmt.Println(child.String())
	fmt.Println("rawdata :", string(child.RawData()))
}
