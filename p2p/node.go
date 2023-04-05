package p2p

import (
	"context"

	"github.com/ipfs/kubo/core"
)

type Node struct {
	node *core.IpfsNode
}

func NewNode() (*Node, error) {
	node, err := core.NewNode(context.Background(), &core.BuildCfg{
		Online: true,
	})
	if err != nil {
		return nil, err
	}
	return &Node{node}, nil
}
