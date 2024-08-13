package config

// 以太坊配置
type ethConfig struct {
	Url             string `json:"Url"`             //RPC连接地址
	ContractAddress string `json:"ContractAddress"` // ERC20合约地址
}
