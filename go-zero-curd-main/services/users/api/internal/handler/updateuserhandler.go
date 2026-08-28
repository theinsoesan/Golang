package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"private/services/users/api/internal/logic"
	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"
)

func updateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUpdateUserLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUser(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
