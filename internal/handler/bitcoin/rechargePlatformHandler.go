package bitcoin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nft/internal/logic/bitcoin"
	"nft/internal/svc"
	"nft/internal/types"
)

func RechargePlatformHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RechargePlatformReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := bitcoin.NewRechargePlatformLogic(r.Context(), svcCtx)
		resp, err := l.RechargePlatform(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
