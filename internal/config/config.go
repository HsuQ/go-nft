package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	BitcoinConf     BitcoinConf       `json:"BitcoinConf"`
	BlockcypherConf BlockcypherConfig `json:"Blockcypher"`
	ETHConf         ethConfig         `json:"ETHConf"`
}
