package bitcoin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nft/internal/logic/bitcoin"
	"nft/internal/svc"
)

func GetBalanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := bitcoin.NewGetBalanceLogic(r.Context(), svcCtx)
		resp, err := l.GetBalance()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
