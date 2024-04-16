package config

// 比特币RPC配置
type BitcoinConf struct {
	Host       string `json:"Host"`       //RPC连接地址
	User       string `json:"User"`       //用户名
	Pass       string `json:"Pass"`       //密码
	WalletName string `json:"WalletName"` //钱包名称
	WalletPass string `json:"WalletPass"` //钱包密码
	Mainnet    bool   `json:"Mainnet"`    //是否是主网
}
