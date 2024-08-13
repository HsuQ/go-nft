package eth

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"nft/internal/svc"
	"nft/internal/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckERC20ArrivedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckERC20ArrivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckERC20ArrivedLogic {
	return &CheckERC20ArrivedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckERC20ArrivedLogic) CheckERC20Arrived(req *types.CheckERC20ArrivedReq) (resp *types.CheckERC20ArrivedResp, err error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(req.Address) {
		logx.Errorf("地址格式不正确")
		return nil, fmt.Errorf("地址格式不正确")
	}
	accountAddress := common.HexToAddress(req.Address)
	contractAddress := common.HexToAddress(l.svcCtx.Config.ETHConf.ContractAddress)

	// Get the latest block number
	blockHeader, err := l.svcCtx.EthClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		logx.Errorf("获取block失败")
		return nil, fmt.Errorf("获取block失败")
	}
	logx.Infof("当前block数量: %v", blockHeader.Number)

	erc20ABI, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"}]`))
	if err != nil {
		logx.Errorf("获取erc20ABI失败")
		return nil, fmt.Errorf("获取erc20ABI失败")
	}
	// 构造balanceOf方法的调用数据
	data, err := erc20ABI.Pack("balanceOf", accountAddress)
	if err != nil {
		logx.Errorf("获取balanceOf方法的调用数据失败")
		return nil, fmt.Errorf("获取balanceOf方法的调用数据失败")
	}
	// data, err := hexutil.Decode("0xb69ef8a8") // First 4 bytes of keccak256("balance()")

	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	response, err := l.svcCtx.EthClient.CallContract(context.Background(), callMsg, blockHeader.Number)

	if err != nil {
		logx.Errorf("获取erc20余额失败： " + err.Error())
		return nil, err
	}
	// 解析返回值
	var balance = new(big.Int)
	balance.SetBytes(response)

	logx.Infof("账号{%v}余额为： %v", req.Address, balance)
	// etherBalance := new(big.Int).Div(balance, big.NewInt(1e18))

	// 将 amountStr 转换为 big.Int
	amountInWei := new(big.Int)
	amountInWei, ok := amountInWei.SetString(req.Amount, 10)
	if !ok {
		logx.Errorf("Invalid amount: %s " + req.Amount)
	}

	resp = &types.CheckERC20ArrivedResp{
		Balance: balance.String(),
		Arrived: balance.Cmp(amountInWei) >= 0,
		Address: req.Address,
	}
	return resp, nil
}
