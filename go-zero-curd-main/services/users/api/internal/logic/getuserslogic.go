package logic

import (
	"context"
	"errors"

	"private/services/users/api/internal/svc"
	"private/services/users/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetusersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetusersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetusersLogic {
	return &GetusersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetusersLogic) Getusers() (resp []*model.User, err error) {

	usersInfo, err := l.svcCtx.UserModel.GetUsers(l.ctx)
	resp = usersInfo
	switch err {
	case nil:
	case model.ErrNotFound:
		return nil, errors.New("There is no data to show")
	default:
		return nil, errors.New("There is no data to show")
	}

	return resp, nil
}
