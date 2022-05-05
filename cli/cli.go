package cli

import (
	"flag"
	"fmt"

	"github.com/figment-networks/oasis-rpc-proxy/client"
	"github.com/figment-networks/oasis-rpc-proxy/config"
	"github.com/figment-networks/oasis-rpc-proxy/utils/logger"
	"github.com/figment-networks/oasis-rpc-proxy/utils/reporting"
	"github.com/oasisprotocol/oasis-core/go/genesis/api"
	genesisFile "github.com/oasisprotocol/oasis-core/go/genesis/file"
	oasisFlags "github.com/oasisprotocol/oasis-core/go/oasis-node/cmd/common/flags"
	"github.com/pkg/errors"
)

// Run executes the command line interface
func Run() {
	defer reporting.RecoverError()

	var configPath string
	var runCommand string
	var showVersion bool

	flag.BoolVar(&showVersion, "v", false, "Show application version")
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.StringVar(&runCommand, "cmd", "", "Command to run")
	flag.Var(
		oasisFlags.DebugDontBlameOasisFlag.Lookup(oasisFlags.CfgDebugDontBlameOasis).Value,
		oasisFlags.CfgDebugDontBlameOasis,
		"Enable debug/unsafe/insecure Oasis options",
	)
	flag.Parse()

	if showVersion {
		fmt.Println(versionString())
		return
	}

	cfg, err := initConfig(configPath)
	if err != nil {
		panic(fmt.Errorf("error initializing config [ERR: %+v]", err))
	}

	if err = initLogger(cfg); err != nil {
		panic(fmt.Errorf("error initializing logger [ERR: %+v]", err))
	}

	initErrorReporting(cfg)

	if runCommand == "" {
		terminate(errors.New("command is required"))
	}

	if err := startCommand(cfg, runCommand); err != nil {
		terminate(err)
	}
}

func startCommand(cfg *config.Config, name string) error {
	switch name {
	case "server":
		return startServer(cfg)
	default:
		return errors.New(fmt.Sprintf("command %s not found", name))
	}
}

func terminate(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func initConfig(path string) (*config.Config, error) {
	cfg := config.New()

	if err := config.FromEnv(cfg); err != nil {
		return nil, err
	}

	if path != "" {
		if err := config.FromFile(path, cfg); err != nil {
			return nil, err
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func initLogger(cfg *config.Config) error {
	_, err := logger.Init(cfg)
	return err
}

func initClient(cfg *config.Config) (*client.Client, error) {
	return client.New(cfg.OasisSocket)
}

func initGenesis(cfg *config.Config) (*api.Document, error) {
	logger.Debug(fmt.Sprintf("genesis file path: %s", cfg.GenesisFilePath))

	genesis, err := genesisFile.NewFileProvider(cfg.GenesisFilePath)
	if err != nil {
		return nil, err
	}

	// Retrieve the genesis document and use it to configure the ChainID for
	// signature domain separation. We do this as early as possible.
	doc, err := genesis.GetGenesisDocument()
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Chain context: '%v'", doc.ChainContext()))

	doc.SetChainContext()

	return doc, nil
}

func initErrorReporting(cfg *config.Config) {
	reporting.Init(cfg)
}
