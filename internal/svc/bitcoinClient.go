package svc

import (
	"fmt"
	"nft/internal/config"

	"github.com/btcsuite/btcd/rpcclient"
)

func NewBitcoinClient(c config.Config) (*rpcclient.Client, error) {
	// 创建RPC客户端连接
	connCfg := &rpcclient.ConnConfig{
		Host:         c.BitcoinConf.Host,
		User:         c.BitcoinConf.User,
		Pass:         c.BitcoinConf.Pass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	btcClient, err := rpcclient.New(connCfg, nil)
	defer btcClient.Shutdown()
	if err != nil {
		fmt.Println("rpcclient.New error " + err.Error())
		return nil, err
	}
	fmt.Println("BTC client connected successfully.")
	return btcClient, nil
}
