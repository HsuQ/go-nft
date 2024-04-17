package config

// blockcypher配置
type BlockcypherConfig struct {
	Token   string `json:"Token"`   //Token
	Coin    string `json:"Coin"`    //币种
	Network string `json:"Network"` //网络
}
