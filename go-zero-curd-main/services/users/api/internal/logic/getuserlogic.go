package logic

import (
	"context"
	"errors"

	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"
	"private/services/users/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetuserLogic {
	return &GetuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetuserLogic) Getuser(req *types.UserReq) (resp *model.User, err error) {
	userInfo, err := l.svcCtx.UserModel.GetUser(l.ctx, req.UserId)
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("There is no data to show")
	default:
		return nil, errors.New("There is no data to show")
	}
	resp = &userInfo
	return resp, nil
}
