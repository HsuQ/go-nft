package svc

import (
	"nft/internal/config"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

func NewEthClient(c config.Config) (*ethclient.Client, error) {
	// 创建以太坊客户端连接
	// client, err := ethclient.Dial("https://cloudflare-eth.com")

	// client, err := ethclient.Dial("https://sepolia.infura.io/v3/69898b33f64f492fb47fcf24349c4291")
	client, err := ethclient.Dial(c.ETHConf.Url)
	if err != nil {
		logx.Errorf("ethclient.Dial error " + err.Error())
		return nil, err
	}
	return client, nil
}
