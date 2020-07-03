package server

import (
	"fmt"
	"net"

	"github.com/figment-networks/oasis-rpc-proxy/client"
	"github.com/figment-networks/oasis-rpc-proxy/config"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/debondingdelegation/debondingdelegationpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/delegation/delegationpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasis-rpc-proxy/utils/logger"
	"github.com/oasisprotocol/oasis-core/go/genesis/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server *grpc.Server

	cfg    *config.Config
	client *client.Client
	doc    *api.Document
}

func New(cfg *config.Config, c *client.Client, doc *api.Document) *Server {
	app := &Server{
		server: grpc.NewServer(),
		cfg:    cfg,
		client: c,
		doc:    doc,
	}
	return app.init()
}

func (s *Server) init() *Server {
	logger.Info("initializing gRPC servers...", logger.Field("app", "server"))

	chainpb.RegisterChainServiceServer(s.server, NewChainServer(s.client, s.doc))
	accountpb.RegisterAccountServiceServer(s.server, NewAccountServer(s.client))
	blockpb.RegisterBlockServiceServer(s.server, NewBlockServer(s.client))
	statepb.RegisterStateServiceServer(s.server, NewStateServer(s.client))
	validatorpb.RegisterValidatorServiceServer(s.server, NewValidatorServer(s.client))
	transactionpb.RegisterTransactionServiceServer(s.server, NewTransactionServer(s.client))
	delegationpb.RegisterDelegationServiceServer(s.server, NewDelegationServer(s.client))
	debondingdelegationpb.RegisterDebondingDelegationServiceServer(s.server, NewDebondingDelegationServer(s.client))

	// Register reflection service on gRPC accountServer.
	reflection.Register(s.server)

	return s
}

func (s *Server) Start(serverAddress string, serverPort int64) error {
	logger.Info(fmt.Sprintf("starting grpc accountServer %s:%d...", serverAddress, serverPort), logger.Field("app", "server"))

	//go s.startMetricsServer()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverAddress, serverPort))
	if err != nil {
		return err
	}
	return s.server.Serve(lis)
}
