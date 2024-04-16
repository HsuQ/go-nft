package svc

import (
	"encoding/json"
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
	if err != nil {
		fmt.Println("rpcclient.New error " + err.Error())
		return nil, err
	}
	// 钱包名称
	walletName := c.BitcoinConf.WalletName

	// 获取所有钱包
	rawResponse, err := btcClient.RawRequest("listwalletdir", []json.RawMessage{})
	if err != nil {
		fmt.Println("ListWalletDir error " + err.Error())
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
		fmt.Println(wallet.Name)
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
			json.RawMessage("true"),
			json.RawMessage(`"` + password + `"`),
		}

		_, err = btcClient.RawRequest("createwallet", params)
		if err != nil {
			fmt.Println("CreateWallet error " + err.Error())
			return nil, err
		}
		fmt.Printf("Wallet %s created successfully.", walletName)
	}

	// 获取已加载的钱包列表
	rawResponse, err = btcClient.RawRequest("listwallets", []json.RawMessage{})
	if err != nil {
		fmt.Println("ListWallets error " + err.Error())
		return nil, err
	}

	var loadedWallets []string
	err = json.Unmarshal(rawResponse, &loadedWallets)
	if err != nil {
		fmt.Println("Error unmarshalling listwallets response " + err.Error())
		return nil, err
	}

	// 检查钱包是否已经被加载
	isLoaded := false
	for _, loadedWallet := range loadedWallets {
		if loadedWallet == walletName {
			isLoaded = true
			fmt.Println("钱包已加载 ")
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
		fmt.Println("Wallet loaded successfully.")
	}
	fmt.Println("BTC client connected successfully.")
	connCfg.Host = c.BitcoinConf.Host + "/wallet/" + walletName

	btcClient, err = rpcclient.New(connCfg, nil)

	if err != nil {
		fmt.Println("rpcclient.New error " + err.Error())
		return nil, err
	}

	// Refill the keypool
	// _, err = btcClient.RawRequest("keypoolrefill", nil)
	// if err != nil {
	// 	return nil, err
	// }

	return btcClient, nil
}
