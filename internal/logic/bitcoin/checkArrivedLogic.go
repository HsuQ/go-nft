package bitcoin

import (
	"context"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckArrivedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// type Address struct {
// 	Address string `json:"address"`
// 	Balance int64  `json:"balance"`
// }

type Address struct {
	Balance int64 `json:"final_balance"`
}

func NewCheckArrivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckArrivedLogic {
	return &CheckArrivedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// func (l *CheckArrivedLogic) CheckArrived(req *types.CheckArrivedReq) (resp *types.CheckArrivedResp, err error) {
// 	address := req.Address
// 	// url := fmt.Sprintf("https://api.blockcypher.com/v1/btc/main/addrs/%s/balance", address)
// 	url := fmt.Sprintf("https://blockchain.info/rawaddr/%s", address)

// 	res, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	var addr Address

// 	// 打印body
// 	// fmt.Println(string(body))

// 	err = json.Unmarshal(body, &addr)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Printf("Balance of address %s: %d\n", req.Address, addr.Balance)

// 	arrived := addr.Balance >= req.Amount

// 	return &types.CheckArrivedResp{
// 		Address: req.Address,
// 		Arrived: arrived,
// 		Balance: addr.Balance,
// 	}, nil
// }

func (l *CheckArrivedLogic) CheckArrived(req *types.CheckArrivedReq) (resp *types.CheckArrivedResp, err error) {
	client := l.svcCtx.BitcoinClient
	address, err := btcutil.DecodeAddress(req.Address, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	// 仅关注确认数6<n<9999999的未消费输出
	utxos, err := client.ListUnspentMinMaxAddresses(6, 9999999, []btcutil.Address{address})
	if err != nil {
		return nil, err
	}
	var total float64
	for _, utxo := range utxos {
		total += utxo.Amount
	}

	arrived := total >= float64(req.Amount)
	return &types.CheckArrivedResp{
		Address: req.Address,
		Arrived: arrived,
		Balance: int64(total),
	}, nil

}
