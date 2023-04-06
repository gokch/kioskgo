package p2p

import (
	"testing"

	"github.com/ipfs/boxo/routing/http/client"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	cli, err := client.New("http://localhost:2379")
	require.NoError(t, err)
	_ = cli
}
