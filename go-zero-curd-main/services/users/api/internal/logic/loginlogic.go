package logic

import (
	"context"
	"private/services/common/errorx"
	"strings"
	"time"

	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"
	"private/services/users/model"

	"github.com/golang-jwt/jwt"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
func (l *LoginLogic) Login(req *types.LoginReq) (resp types.LoginReply, err error) {

	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return resp, errorx.NewDefaultError("Invalid parameter")
	}

	userInfo, err := l.svcCtx.UserModel.Login(l.ctx, req.Username, req.Password)

	switch err {
	case nil:
	case model.ErrNotFound:
		return resp, errorx.NewDefaultError("Username does not exist")
	default:
		return resp, err
	}

	if userInfo.Password != req.Password {
		return resp, errorx.NewDefaultError("User password is incorrect")
	}

	// ---start---
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, userInfo.Id)
	if err != nil {
		return resp, err
	}
	// ---end---

	return types.LoginReply{
		Id:           userInfo.Id,
		Name:         userInfo.Name,
		Gender:       userInfo.Gender,
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}
