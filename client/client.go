package client

import (
	"fmt"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/oasis-rpc-proxy/utils/logger"
	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	oasisGrpc "github.com/oasisprotocol/oasis-core/go/common/grpc"
	"github.com/oasisprotocol/oasis-core/go/staking/api"

	"google.golang.org/grpc"
)

var clientRequestDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
	Namespace: "indexers",
	Subsystem: "oasis_proxy",
	Name:      "node_request_duration",
	Desc:      "The total time required to execute request to node",
	Tags:      []string{"request"},
})

func New(target string) (*Client, error) {
	logger.Debug(fmt.Sprintf("grpc server target is %s", target))

	conn, err := oasisGrpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,

		Consensus: NewConsensusClient(conn),
		Registry:  NewRegistryClient(conn),
		Scheduler: NewSchedulerClient(conn),
		Staking:   NewStakingClient(conn),
	}, nil
}

type Client struct {
	conn *grpc.ClientConn

	Consensus ConsensusClient
	Registry  RegistryClient
	Scheduler SchedulerClient
	Staking   StakingClient
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func getPublicKey(key string) (*signature.PublicKey, error) {
	var pKey signature.PublicKey
	if err := pKey.UnmarshalText([]byte(key)); err != nil {
		return nil, err
	}
	return &pKey, nil
}

func getAddress(rawAddress string) (*api.Address, error) {
	address := api.Address{}
	if err := address.UnmarshalText([]byte(rawAddress)); err != nil {
		return nil, err
	}
	return &address, nil
}

func getAddressFromPublicKey(key string) (*api.Address, error) {
	var pk signature.PublicKey
	if err := pk.UnmarshalText([]byte(key)); err != nil {
		return nil, err
	}
	address := api.NewAddress(pk)
	return &address, nil
}
