package bitcoin

import (
	"context"
	"fmt"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
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

func (l *NewAccountLogic) NewAccount() (resp *types.NewAccountResp, err error) {

	// 生成新的密钥对
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}

	pubKey := privKey.PubKey()

	// 获取关联的比特币地址
	address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取私钥的WIF格式
	wif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Generated Address: " + address.EncodeAddress())
	fmt.Println("Private Key WIF: " + wif.String())

	return &types.NewAccountResp{
		Address:    address.EncodeAddress(),
		PrivateKey: wif.String(),
	}, nil
}
