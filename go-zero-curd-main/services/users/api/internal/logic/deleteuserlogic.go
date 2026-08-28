package logic

import (
	"context"
	"errors"
	"private/services/users/model"

	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.UserDeleteReq) (resp types.UserDeleteReply, err error) {

	userInfo, err := l.svcCtx.UserModel.DeleteUser(l.ctx, req.Id)
	switch err {
	case nil:
	case model.ErrNotFound:
		return resp, errors.New("There is no data to show")
	default:
		return resp, errors.New("There is no data to show")
	}
	return types.UserDeleteReply{
		ResCode: userInfo.ResCode,
		ResDesc: userInfo.ResDesc,
	}, nil
}
