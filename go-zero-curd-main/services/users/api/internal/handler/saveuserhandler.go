package handler

import (
	"net/http"

	"private/services/users/api/internal/logic"
	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func saveuserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSaveuserLogic(r.Context(), svcCtx)
		resp, err := l.Saveuser(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
