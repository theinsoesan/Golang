package logic

import (
	"context"
	"errors"
	"private/services/users/model"

	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UserUpdateReq) (resp types.UserUpdateReply, err error) {
	var reqUser model.User
	reqUser.Id = req.Id
	reqUser.Name = req.Name
	reqUser.Password = req.Password
	reqUser.Gender = req.Gender
	userInfo, err := l.svcCtx.UserModel.UpdateUser(l.ctx, reqUser)
	switch err {
	case nil:
	case model.ErrNotFound:
		return resp, errors.New("There is no data to show")
	default:
		return resp, errors.New("There is no data to show")
	}
	return types.UserUpdateReply{
		ResCode: userInfo.ResCode,
		ResDesc: userInfo.ResDesc,
	}, nil
}
