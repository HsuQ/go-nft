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

type GetTransactionCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransactionCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionCountLogic {
	return &GetTransactionCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransactionCountLogic) GetTransactionCount(req *types.GetTransactionCountReq) (resp *types.GetTransactionCountResp, err error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(req.Address) {
		logx.Errorf("地址格式不正确")
		return nil, fmt.Errorf("地址格式不正确")
	}

	account := common.HexToAddress(req.Address)

	// 获取 nonce
	nonce, err := l.svcCtx.EthClient.PendingNonceAt(l.ctx, account)
	if err != nil {
		logx.Errorf("获取eth nonce失败： " + err.Error())
		return nil, err
	}

	// 获取当前推荐的 gas 价格
	gasPrice, err := l.svcCtx.EthClient.SuggestGasPrice(l.ctx)
	if err != nil {
		logx.Errorf("获取gas price失败： " + err.Error())
		return nil, err
	}

	// 获取 maxPriorityFeePerGas (EIP-1559)
	maxPriorityFeePerGas, err := l.svcCtx.EthClient.SuggestGasTipCap(l.ctx)
	if err != nil {
		logx.Errorf("获取maxPriorityFeePerGas失败： " + err.Error())
		return nil, err
	}

	logx.Infof("账号{%v} nonce为： %v, gasPrice: %v wei, maxPriorityFeePerGas: %v wei",
		req.Address, nonce, gasPrice, maxPriorityFeePerGas)

	resp = &types.GetTransactionCountResp{
		Nonce:                strconv.FormatUint(nonce, 10),
		Address:              req.Address,
		GasPrice:             gasPrice.String(),
		MaxPriorityFeePerGas: maxPriorityFeePerGas.String(),
	}

	return resp, nil
}
