package bitcoin

import (
	"context"

	"nft/internal/svc"
	"nft/internal/types"

	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
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
	// 创建RPC客户端连接
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "yourusername",
		Pass:         "yourpassword",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// 创建比特币地址
	address, err := btcutil.DecodeAddress("yourBitcoinAddress", &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	// 获取地址的未消费输出
	unspent, err := client.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{address})
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

	// 发送转账请求
	_, err = client.SendToAddress(address, amount)
	if err != nil {
		log.Fatal(err)
		return &types.RechargePlatformResp{
			Code: 500,
			Msg:  "Transaction failed",
		}, err
	}

	log.Println("Transaction successful")
	return &types.RechargePlatformResp{
		Code: 200,
		Msg:  "Transaction successful",
	}, nil
}
