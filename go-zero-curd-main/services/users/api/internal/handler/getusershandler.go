package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"private/services/users/api/internal/logic"
	"private/services/users/api/internal/svc"
)

func getusersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetusersLogic(r.Context(), svcCtx)
		resp, err := l.Getusers()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
