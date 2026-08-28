package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"private/services/users/api/internal/logic"
	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"
)

func deleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDeleteUserLogic(r.Context(), svcCtx)
		resp, err := l.DeleteUser(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
