package bitcoin

import (
	"net/http"

	"nft/internal/logic/bitcoin"
	"nft/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func NewAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := bitcoin.NewNewAccountLogic(r.Context(), svcCtx)
		resp, err := l.NewWalletAccount()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
