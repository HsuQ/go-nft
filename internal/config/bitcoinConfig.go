package config

// 比特币RPC配置
type BitcoinConf struct {
	Host string `json:"Host"` //RPC连接地址
	User string `json:"User"` //用户名
	Pass string `json:"Pass"` //密码
}
