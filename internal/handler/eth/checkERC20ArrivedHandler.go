package eth

import (
	"net/http"

	"nft/internal/logic/eth"
	"nft/internal/svc"
	"nft/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckERC20ArrivedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckERC20ArrivedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := eth.NewCheckERC20ArrivedLogic(r.Context(), svcCtx)
		resp, err := l.CheckERC20Arrived(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
