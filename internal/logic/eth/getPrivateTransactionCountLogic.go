package eth

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPrivateTransactionCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPrivateTransactionCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPrivateTransactionCountLogic {
	return &GetPrivateTransactionCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPrivateTransactionCountLogic) GetPrivateTransactionCount(req *types.GetTransactionCountReq) (resp *types.GetTransactionCountResp, err error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(req.Address) {
		logx.Errorf("地址格式不正确")
		return nil, fmt.Errorf("地址格式不正确")
	}
	account := common.HexToAddress(req.Address)
	nonce, err := l.svcCtx.PrivateEthClient.PendingNonceAt(l.ctx, account)
	if err != nil {
		logx.Errorf("获取eth nonce失败： " + err.Error())
		return nil, err
	}

	logx.Infof("账号{%v} nonce为： %v", req.Address, nonce)
	nonceStr := strconv.FormatUint(nonce, 10) // 转换为十进制字符串
	resp = &types.GetTransactionCountResp{
		Nonce:   nonceStr,
		Address: req.Address,
	}

	return resp, nil
}
