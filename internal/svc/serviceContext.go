package svc

import (
	"nft/internal/config"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ServiceContext struct {
	Config        config.Config
	BitcoinClient *rpcclient.Client
	EthClient     *ethclient.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	btcClient, err := NewBitcoinClient(c)
	if err != nil {
		// handle error
		panic(err)
	}
	ethClient, err := NewEthClient(c)
	if err != nil {
		// handle error
		panic(err)
	}
	return &ServiceContext{
		Config:        c,
		BitcoinClient: btcClient,
		EthClient:     ethClient,
	}
}
