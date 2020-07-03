package client

import (
	"context"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/oasisprotocol/oasis-core/go/common/node"
	"github.com/oasisprotocol/oasis-core/go/registry/api"
	"google.golang.org/grpc"
)

var (
	_ RegistryClient = (*registryClient)(nil)
)

type RegistryClient interface {
	GeNodeById(context.Context, string, int64) (*node.Node, error)
}

func NewRegistryClient(conn *grpc.ClientConn) *registryClient {
	return &registryClient{
		client: api.NewRegistryClient(conn),
	}
}

type registryClient struct {
	client api.Backend
}

func (c *registryClient) GeNodeById(ctx context.Context, key string, height int64) (*node.Node, error) {
	pKey, err := getPublicKey(key)
	if err != nil {
		return nil, err
	}
	q := &api.IDQuery{
		Height: height,
		ID:     *pKey,
	}

	t := metrics.NewTimer(clientRequestDuration.WithLabels([]string{"RegistryClient_GeNodeById"}))
	defer t.ObserveDuration()
	return c.client.GetNode(ctx, q)
}
