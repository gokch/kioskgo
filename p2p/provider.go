package p2p

import (
	"context"
	"time"

	"github.com/ipfs/boxo/routing/http/server"
	"github.com/ipfs/boxo/routing/http/types"
	"github.com/ipfs/go-cid"
)

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) FindProviders(ctx context.Context, key cid.Cid) []types.ProviderResponse {
	return nil
}

func (p *Provider) ProvideBitswap(ctx context.Context, req *server.BitswapWriteProvideRequest) (time.Duration, error) {
	return time.Second, nil
}

func (p *Provider) Provide(ctx context.Context, req *server.WriteProvideRequest) (types.ProviderResponse, error) {
	return nil, nil
}
