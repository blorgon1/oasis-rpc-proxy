package client

import (
	"context"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/oasisprotocol/oasis-core/go/scheduler/api"
	"google.golang.org/grpc"
)

var (
	_ SchedulerClient = (*schedulerClient)(nil)
)

type SchedulerClient interface {
	GetValidatorsByHeight(context.Context, int64) ([]*api.Validator, error)
}

func NewSchedulerClient(conn *grpc.ClientConn) SchedulerClient {
	return &schedulerClient{
		client: api.NewSchedulerClient(conn),
	}
}

type schedulerClient struct {
	client api.Backend
}

func (r *schedulerClient) GetValidatorsByHeight(ctx context.Context, h int64) ([]*api.Validator, error) {
	t := metrics.NewTimer(clientRequestDuration.WithLabels([]string{"SchedulerClient_GetValidatorsByHeight"}))
	defer t.ObserveDuration()
	return r.client.GetValidators(ctx, h)
}
