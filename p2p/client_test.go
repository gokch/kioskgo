package p2p

import (
	"context"
	"fmt"
	"testing"

	"github.com/ipfs/go-libipfs/blocks"
)

func TestIpfs(t *testing.T) {
	context := context.Background()

	store1, err := NewIPFSstore(context, true)
	if err != nil {
		t.Fatal(err)
	}

	blockData := blocks.NewBlock([]byte("hello world!!"))
	err = store1.Put(context, blockData)
	if err != nil {
		t.Fatal(err)
	}

	store2, err := NewIPFSstore(context, true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("cid : %v\n", blockData.Cid())

	blockNew, err := store2.Get(context, blockData.Cid())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(blockNew.RawData())
}
