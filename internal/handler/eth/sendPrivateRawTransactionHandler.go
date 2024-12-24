package eth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nft/internal/logic/eth"
	"nft/internal/svc"
	"nft/internal/types"
)

func SendPrivateRawTransactionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendRawTransactionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := eth.NewSendPrivateRawTransactionLogic(r.Context(), svcCtx)
		resp, err := l.SendPrivateRawTransaction(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
