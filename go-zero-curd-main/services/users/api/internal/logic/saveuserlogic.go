package logic

import (
	"context"
	"errors"

	"private/services/users/api/internal/svc"
	"private/services/users/api/internal/types"
	"private/services/users/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveuserLogic {
	return &SaveuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveuserLogic) Saveuser(req types.UserSaveReq) (resp types.UserSaveReply, err error) {
	var reqUser model.User
	reqUser.Name = req.Name
	reqUser.Password = req.Password
	reqUser.Gender = req.Gender
	userInfo, err := l.svcCtx.UserModel.SaveUser(l.ctx, reqUser)
	switch err {
	case nil:
	case model.ErrNotFound:
		return resp, errors.New("There is no data to show")
	default:
		return resp, errors.New("There is no data to show")
	}
	return types.UserSaveReply{
		ResCode: userInfo.ResCode,
		ResDesc: userInfo.ResDesc,
	}, nil
}
