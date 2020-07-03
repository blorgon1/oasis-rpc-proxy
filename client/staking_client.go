package client

import (
	"context"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/oasisprotocol/oasis-core/go/staking/api"
	"google.golang.org/grpc"
)

var (
	_ StakingClient = (*stakingClient)(nil)
)

type StakingClient interface {
	GetAccountByAddress(context.Context, string, int64) (*api.Account, error)
	GetDelegations(context.Context, string, int64) (map[api.Address]*api.Delegation, error)
	GetDebondingDelegations(context.Context, string, int64) (map[api.Address][]*api.DebondingDelegation, error)
	GetState(context.Context, int64) (*api.Genesis, error)
}

func NewStakingClient(conn *grpc.ClientConn) *stakingClient {
	return &stakingClient{
		client: api.NewStakingClient(conn),
	}
}

type stakingClient struct {
	client api.Backend
}

func (c *stakingClient) GetAccountByAddress(ctx context.Context, key string, height int64) (*api.Account, error) {
	q, err := c.buildOwnerQuery(key, height)
	if err != nil {
		return nil, err
	}

	t := metrics.NewTimer(clientRequestDuration.WithLabels([]string{"StakingClient_GetAccountByAddress"}))
	defer t.ObserveDuration()
	return c.client.Account(ctx, q)
}

func (c *stakingClient) GetDelegations(ctx context.Context, key string, height int64) (map[api.Address]*api.Delegation, error) {
	q, err := c.buildOwnerQuery(key, height)
	if err != nil {
		return nil, err
	}

	t := metrics.NewTimer(clientRequestDuration.WithLabels([]string{"StakingClient_GetDelegations"}))
	defer t.ObserveDuration()
	return c.client.Delegations(ctx, q)
}

func (c *stakingClient) GetDebondingDelegations(ctx context.Context, key string, height int64) (map[api.Address][]*api.DebondingDelegation, error) {
	q, err := c.buildOwnerQuery(key, height)
	if err != nil {
		return nil, err
	}

	t := metrics.NewTimer(clientRequestDuration.WithLabels([]string{"StakingClient_GetDebondingDelegations"}))
	defer t.ObserveDuration()
	return c.client.DebondingDelegations(ctx, q)
}

func (c *stakingClient) GetState(ctx context.Context, height int64) (*api.Genesis, error) {
	t := metrics.NewTimer(clientRequestDuration.WithLabels([]string{"StakingClient_GetState"}))
	defer t.ObserveDuration()
	return c.client.StateToGenesis(ctx, height)
}

func (c *stakingClient) buildOwnerQuery(key string, height int64) (*api.OwnerQuery, error) {
	address, err := getAddress(key)
	if err != nil {
		return nil, err
	}
	q := &api.OwnerQuery{
		Height: height,
		Owner:  *address,
	}
	return q, nil
}
