package bitcoin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckArrivedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

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

func (l *CheckArrivedLogic) CheckArrived(req *types.CheckArrivedReq) (resp *types.CheckArrivedResp, err error) {
	if l.svcCtx.Config.BitcoinConf.Mainnet {
		return l.CheckArrivedMainNet(req)
	} else {
		return l.CheckArrivedTestnet3(req)
	}
}

func (l *CheckArrivedLogic) CheckArrivedMainNet(req *types.CheckArrivedReq) (resp *types.CheckArrivedResp, err error) {
	address := req.Address
	// url := fmt.Sprintf("https://api.blockcypher.com/v1/btc/main/addrs/%s/balance", address)
	url := fmt.Sprintf("https://blockchain.info/rawaddr/%s", address)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var addr Address

	// 打印body
	// fmt.Println(string(body))

	err = json.Unmarshal(body, &addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Balance of address %s: %d\n", req.Address, addr.Balance)

	arrived := addr.Balance >= req.Amount

	return &types.CheckArrivedResp{
		Address: req.Address,
		Arrived: arrived,
		Balance: addr.Balance,
	}, nil
}

type AddressBalance struct {
	Balance int64 `json:"balance"`
}

func (l *CheckArrivedLogic) CheckArrivedTestnet3(req *types.CheckArrivedReq) (myResp *types.CheckArrivedResp, err error) {
	address := req.Address
	url := fmt.Sprintf("https://api.blockcypher.com/v1/btc/test3/addrs/%s/balance", address)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var balance AddressBalance
	err = json.Unmarshal(body, &balance)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Balance: %d satoshis\n", balance.Balance)

	arrived := balance.Balance >= req.Amount

	return &types.CheckArrivedResp{
		Address: req.Address,
		Arrived: arrived,
		Balance: balance.Balance,
	}, nil
}
