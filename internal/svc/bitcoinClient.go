package svc

import (
	"encoding/json"
	"fmt"
	"nft/internal/config"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/zeromicro/go-zero/core/logx"
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
	if err != nil {
		logx.Errorf("rpcclient.New error " + err.Error())
		return nil, err
	}
	// 钱包名称
	walletName := c.BitcoinConf.WalletName

	// 获取所有钱包
	rawResponse, err := btcClient.RawRequest("listwalletdir", []json.RawMessage{})
	if err != nil {
		logx.Infof("ListWalletDir error " + err.Error())
		return nil, err
	}

	var walletDir struct {
		Wallets []struct {
			Name string `json:"name"`
		} `json:"wallets"`
	}
	err = json.Unmarshal(rawResponse, &walletDir)
	if err != nil {
		fmt.Println("Error unmarshalling listwalletdir response " + err.Error())
		return nil, err
	}

	// 检查钱包是否存在
	walletExists := false
	for _, wallet := range walletDir.Wallets {
		if wallet.Name == walletName {
			walletExists = true
			break
		}
	}

	// 如果钱包不存在，则创建钱包
	if !walletExists {
		password := c.BitcoinConf.WalletPass
		params := []json.RawMessage{
			json.RawMessage(`"` + walletName + `"`),
			json.RawMessage("false"),
			json.RawMessage("true"), // (boolean, optional, default=false) Create a blank wallet. A blank wallet has no keys or HD seed. One can be set using sethdseed.
			json.RawMessage(`"` + password + `"`),
		}

		_, err = btcClient.RawRequest("createwallet", params)
		if err != nil {
			logx.Errorf("CreateWallet error " + err.Error())
			return nil, err
		}
		logx.Infof("Wallet %s created successfully.", walletName)
	}

	// 获取已加载的钱包列表
	rawResponse, err = btcClient.RawRequest("listwallets", []json.RawMessage{})
	if err != nil {
		logx.Errorf("ListWallets error " + err.Error())
		return nil, err
	}

	var loadedWallets []string
	err = json.Unmarshal(rawResponse, &loadedWallets)
	if err != nil {
		logx.Errorf("Error unmarshalling listwallets response " + err.Error())
		return nil, err
	}

	// 检查钱包是否已经被加载
	isLoaded := false
	for _, loadedWallet := range loadedWallets {
		if loadedWallet == walletName {
			isLoaded = true
			logx.Info("钱包已加载")
			break
		}
	}

	// 如果钱包没有被加载，则加载钱包
	if !isLoaded {
		_, err = btcClient.RawRequest("loadwallet", []json.RawMessage{json.RawMessage(`"` + walletName + `"`)})
		if err != nil {
			fmt.Println("LoadWallet error " + err.Error())
			return nil, err
		}
		logx.Infof("Wallet %s loaded successfully.", walletName)
	}
	logx.Info("BTC client connected successfully.")
	connCfg.Host = c.BitcoinConf.Host + "/wallet/" + walletName

	btcClient, err = rpcclient.New(connCfg, nil)

	if err != nil {
		fmt.Println("rpcclient.New error " + err.Error())
		return nil, err
	}

	// Backup the seed
	backupFile := fmt.Sprintf("/home/%s.txt", walletName)
	// backupFile := path.Join(".", walletName+".txt")
	logx.Infof("Backing up wallet %s to file %s.", walletName, backupFile)
	dumpParams := []json.RawMessage{
		json.RawMessage(`"` + backupFile + `"`),
	}
	_, err = btcClient.RawRequest("dumpwallet", dumpParams)
	if err != nil {
		logx.Errorf("DumpWallet error " + err.Error())
		return nil, err
	}
	logx.Infof("Wallet %s has been backed up to file %s.", walletName, backupFile)

	return btcClient, nil
}
