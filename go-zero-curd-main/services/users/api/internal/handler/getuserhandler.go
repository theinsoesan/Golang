package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"private/services/users/api/internal/logic"
	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"
)

func getuserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetuserLogic(r.Context(), svcCtx)
		resp, err := l.Getuser(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
