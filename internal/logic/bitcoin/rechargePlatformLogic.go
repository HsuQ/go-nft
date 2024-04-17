package bitcoin

import (
	"context"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/zeromicro/go-zero/core/logx"
)

type RechargePlatformLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRechargePlatformLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RechargePlatformLogic {
	return &RechargePlatformLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RechargePlatformLogic) RechargePlatform(req *types.RechargePlatformReq) (resp *types.RechargePlatformResp, err error) {
	client := l.svcCtx.BitcoinClient
	// 将to地址转换为btcutil.Address类型
	var params *chaincfg.Params
	if l.svcCtx.Config.BitcoinConf.Mainnet {
		params = &chaincfg.MainNetParams
	} else {
		params = &chaincfg.TestNet3Params
	}
	toAddress, err := btcutil.DecodeAddress(req.ToAddress, params)
	if err != nil {
		logx.Errorf("error decoding address: %v", err)
		return nil, err
	}

	// balance, err := client.GetBalance("*")
	// if err != nil {
	// 	logx.Error("client.GetBalance error " + err.Error())
	// 	return nil, err
	// }
	balance := btcutil.Amount(600)

	// 解锁钱包60s
	err = client.WalletPassphrase(l.svcCtx.Config.BitcoinConf.WalletPass, 60)
	if err != nil {
		logx.Errorf("error unlocking wallet: %v", err)
	}

	// send the balance to the new address
	txHash, err := client.SendToAddress(toAddress, balance)
	if err != nil {
		logx.Errorf("error sending to address: %v", err)
	}

	logx.Infof("Transaction successful with transaction hash: %s", txHash.String())
	return &types.RechargePlatformResp{
		Code: 200,
		Msg:  "Transaction successful, transaction hash: " + txHash.String(),
	}, nil
}
