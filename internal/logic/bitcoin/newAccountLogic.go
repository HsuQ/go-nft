package bitcoin

import (
	"context"
	"fmt"
	"log"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type NewAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNewAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewAccountLogic {
	return &NewAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NewAccountLogic) NewWalletAccount() (resp *types.NewAccountResp, err error) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "Secret Passphrase")

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	// Display mnemonic and keys
	fmt.Println("Mnemonic: ", mnemonic)
	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)
	var params *chaincfg.Params
	if l.svcCtx.Config.BitcoinConf.Mainnet {
		params = &chaincfg.MainNetParams
	} else {
		params = &chaincfg.TestNet3Params
	}

	address, err := btcutil.NewAddressPubKey(publicKey.Key, params)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Address: ", address)

	return &types.NewAccountResp{
		Address:    address.EncodeAddress(),
		PrivateKey: masterKey.String(),
	}, nil

}

func (l *NewAccountLogic) NewHDWalletAccount() (resp *types.NewAccountResp, err error) {
	client := l.svcCtx.BitcoinClient
	address, err := client.GetNewAddress("")
	if err != nil {
		fmt.Println("client.GetNewAddress error " + err.Error())
		return nil, err
	}
	return &types.NewAccountResp{
		Address:    address.EncodeAddress(),
		PrivateKey: "123",
	}, nil
}
