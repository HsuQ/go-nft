package eth

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeromicro/go-zero/core/logx"
)

type VerifySignatureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifySignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifySignatureLogic {
	return &VerifySignatureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifySignatureLogic) VerifySignature(req *types.VerifySignatureReq) (resp *types.VerifySignatureResp, err error) {
	// 移除可能的 "0x" 前缀
	req.Signature = strings.TrimPrefix(req.Signature, "0x")

	// 签名长度必须是 65 字节 (32 bytes r + 32 bytes s + 1 byte v)
	if len(req.Signature) != 130 {
		logx.Errorf("Invalid signature length: %d", len(req.Signature))
		return nil, errors.New("invalid signature length")
	}

	// 将签名从十六进制解析为字节数组
	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		logx.Errorf("Failed to decode signature: %v", err)
		return nil, errors.New("failed to decode signature")
	}

	// 分割 r, s, v
	r := signature[:32]
	s := signature[32:64]
	v := signature[64]

	// 将 v 从 27/28 转换为 0/1 (以太坊签名规范)
	if v < 27 {
		logx.Errorf("Invalid recovery id (v): %d", v)
		return nil, errors.New("invalid recovery id")
	}
	v -= 27

	// 计算以太坊签名消息的哈希
	prefix := fmt.Sprintf("\u0019Ethereum Signed Message:\n%d", len(req.Message))
	prefixedMessage := []byte(prefix + req.Message)
	messageHash := crypto.Keccak256(prefixedMessage)

	// 组装完整签名 (r, s, v)
	signatureWithV := append(append(r, s...), v)

	// 从签名中恢复公钥
	publicKey, err := crypto.Ecrecover(messageHash, signatureWithV)
	if err != nil {
		logx.Errorf("Failed to recover public key: %v", err)
		return nil, errors.New("failed to recover public key")
	}

	// 将公钥转换为地址
	//recoveredAddress := crypto.PubkeyToAddress(*crypto.UnmarshalPubkey(publicKey)).Hex()
	// 将公钥转换为 ECDSA 格式
	pubKeyECDSA, err := crypto.UnmarshalPubkey(publicKey)
	if err != nil {
		logx.Errorf("Failed to unmarshal public key: %v", err)
		return nil, errors.New("failed to unmarshal public key")
	}

	// 生成地址
	recoveredAddress := crypto.PubkeyToAddress(*pubKeyECDSA).Hex()

	// 对比地址
	return &types.VerifySignatureResp{
		Valid: strings.EqualFold(recoveredAddress, req.Address),
	}, nil
}
