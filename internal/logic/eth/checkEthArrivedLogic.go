package eth

import (
	"context"
	"fmt"
	"regexp"

	"nft/internal/svc"
	"nft/internal/types"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckEthArrivedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckEthArrivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckEthArrivedLogic {
	return &CheckEthArrivedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckEthArrivedLogic) CheckEthArrived(req *types.CheckEthArrivedReq) (resp *types.CheckEthArrivedResp, err error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(req.Address) {
		logx.Errorf("地址格式不正确")
		return nil, fmt.Errorf("地址格式不正确")
	}
	account := common.HexToAddress(req.Address)
	balance, err := l.svcCtx.EthClient.BalanceAt(l.ctx, account, nil)
	if err != nil {
		logx.Errorf("获取eth余额失败： " + err.Error())
		return nil, err
	}
	logx.Infof("账号{%v}余额为： %v", req.Address, balance)
	resp = &types.CheckEthArrivedResp{
		Balance: balance.Int64(),
		Arrived: balance.Cmp(big.NewInt(req.Amount)) >= 0,
		Address: req.Address,
	}
	return resp, nil
}
