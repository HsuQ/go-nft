package eth

import (
	"context"
	"encoding/hex"
	"errors"

	"nft/internal/svc"
	"nft/internal/types"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendRawTransactionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendRawTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRawTransactionLogic {
	return &SendRawTransactionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendRawTransactionLogic) SendRawTransaction(req *types.SendRawTransactionReq) (resp *types.SendRawTransactionResp, err error) {
	// 检查 rawTx 是否为空
	if len(req.RawTx) == 0 {
		return nil, errors.New("rawTx cannot be empty")
	}

	// 将 rawTx（十六进制字符串）解码为字节数组
	rawTxBytes, err := hex.DecodeString(req.RawTx)
	if err != nil {
		return nil, err
	}

	// 解码成 Transaction 格式
	tx := &ethTypes.Transaction{}
	err = rlp.DecodeBytes(rawTxBytes, tx)
	if err != nil {
		return nil, err
	}

	// 调用 SendTransaction 发送交易
	err = l.svcCtx.EthClient.SendTransaction(l.ctx, tx)
	if err != nil {
		return nil, err
	}

	// 返回响应结果
	resp = &types.SendRawTransactionResp{
		TxHash: tx.Hash().Hex(), // 返回交易哈希
	}

	return resp, nil
}
