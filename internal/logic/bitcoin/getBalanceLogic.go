package bitcoin

import (
	"context"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBalanceLogic {
	return &GetBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBalanceLogic) GetBalance() (resp *types.BalanceResp, err error) {
	client := l.svcCtx.BitcoinClient
	balance, err := client.GetBalance("*")
	if err != nil {
		logx.Error("client.GetBalance error " + err.Error())
		return nil, err
	}

	logx.Infof("Balance: %s\n", balance.String())

	return &types.BalanceResp{
		Balance: int64(balance),
	}, nil
}
