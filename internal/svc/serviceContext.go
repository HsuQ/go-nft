package svc

import (
	"nft/internal/config"

	"github.com/btcsuite/btcd/rpcclient"
)

type ServiceContext struct {
	Config        config.Config
	BitcoinClient *rpcclient.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	btcClient, err := NewBitcoinClient(c)
	if err != nil {
		// handle error
		panic(err)
	}
	return &ServiceContext{
		Config:        c,
		BitcoinClient: btcClient,
	}
}
