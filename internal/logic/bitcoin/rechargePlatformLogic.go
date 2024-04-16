package bitcoin

import (
	"context"

	"nft/internal/svc"
	"nft/internal/types"

	"log"

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
	toAddress, err := btcutil.DecodeAddress(req.To, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	// 获取地址的未消费输出
	unspent, err := client.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{toAddress})
	if err != nil {
		log.Fatal(err)
	}

	// 计算总余额
	var total btcutil.Amount
	for _, u := range unspent {
		amount, err := btcutil.NewAmount(u.Amount)
		if err != nil {
			log.Fatal(err)
		}
		total += amount
	}

	// 估算交易的大小和矿工费用
	// 这里假设每个输入的大小为 148 字节，每个输出的大小为 34 字节
	// 并且矿工费用为 1 satoshi/byte
	txSize := len(unspent)*148 + 34 + 10
	fee := btcutil.Amount(txSize)

	// 计算可以转出的金额
	amount := total - fee
	log.Default().Println("Total balance:", amount)

	// 发送转账请求
	// _, err = client.SendToAddress(address, amount)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return &types.RechargePlatformResp{
	// 		Code: 500,
	// 		Msg:  "Transaction failed",
	// 	}, err
	// }

	log.Println("Transaction successful")
	return &types.RechargePlatformResp{
		Code: 200,
		Msg:  "Transaction successful",
	}, nil
}

// func (l *RechargePlatformLogic) RechargePlatformClient(req *types.RechargePlatformReq) (resp *types.RechargePlatformResp, err error) {
// 	client := l.svcCtx.BitcoinClient
// 	amount := 0.01      // 转账金额，单位为 BTC
// 	from := req.From    // 转账账户
// 	toAddress := req.To // 目标地址
// 	toAddressPubKeyHash, err := btcutil.DecodeAddress(toAddress, &chaincfg.MainNetParams)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	txHash, err := client.SendFrom(from, toAddressPubKeyHash, btcutil.Amount(amount*1e8)) // 注意：btcutil.Amount的单位是Satoshi，1 BTC = 1e8 Satoshi
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("Transaction ID: %s\n", txHash)
// }
