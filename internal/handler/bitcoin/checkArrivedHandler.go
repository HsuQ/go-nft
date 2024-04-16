package bitcoin

import (
	"net/http"

	"nft/internal/logic/bitcoin"
	"nft/internal/svc"
	"nft/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckArrivedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckArrivedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := bitcoin.NewCheckArrivedLogic(r.Context(), svcCtx)
		resp, err := l.CheckArrived(&req)
		//
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
