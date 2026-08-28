package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		Login(ctx context.Context, username string, password string) (User, error)
		GetUsers(ctx context.Context) ([]*User, error)
		GetUser(ctx context.Context, userId int64) (User, error)
		SaveUser(ctx context.Context, request User) (resData, error)
		UpdateUser(ctx context.Context, request User) (resData, error)
		DeleteUser(ctx context.Context, request int64) (resData, error)
	}

	customUserModel struct {
		conn sqlx.SqlConn

		*defaultUserModel
	}
	resData struct {
		ResCode string `json:"rescode"`
		ResDesc string `json:"resdesc"`
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &customUserModel{
		conn:             conn,
		defaultUserModel: newUserModel(conn, c),
	}
}
func (m *customUserModel) Login(ctx context.Context, username string, password string) (User, error) {
	sb := squirrel.Select(userRows).From(m.table)
	sb = sb.Where("`name`=?", username)
	query, values, err := sb.ToSql()
	var resp User
	if err := m.conn.QueryRowCtx(ctx, &resp, query, values...); err != nil {
		fmt.Println(resp)

		return resp, err
	}
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return resp, ErrNotFound
	default:
		return resp, err
	}
}
func (m *customUserModel) GetUsers(ctx context.Context) ([]*User, error) {
	sb := squirrel.Select(userRows).From(m.table)
	query, values, err := sb.ToSql()
	resUsers := make([]*User, 0)
	if err := m.conn.QueryRowsCtx(ctx, &resUsers, query, values...); err != nil {
		return nil, err
	}
	switch err {
	case nil:
		return resUsers, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *customUserModel) GetUser(ctx context.Context, userId int64) (User, error) {
	sb := squirrel.Select(userRows).From(m.table)
	sb = sb.Where("`id`=?", userId)
	query, values, err := sb.ToSql()
	var resp User
	if err := m.conn.QueryRowCtx(ctx, &resp, query, values...); err != nil {
		return resp, err
	}
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return resp, ErrNotFound
	default:
		return resp, err
	}
}
func (m *customUserModel) UpdateUser(ctx context.Context, request User) (resData, error) {
	fmt.Println(request)
	var res resData
	fmt.Println("=============> Update")
	data, err := m.FindOne(ctx, request.Id)
	if data != nil {
		if err != nil {
			return res, err
		}
		userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
		//userNumberKey := fmt.Sprintf("%s%v", cacheUserNumberPrefix, data.Number)
		_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
			query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
			fmt.Println("-------------=======>", query)
			return conn.ExecCtx(ctx, query, request.Name, request.Password, request.Gender, request.Id)
		}, userIdKey)
	}
	return res, nil
}
func (m *customUserModel) SaveUser(ctx context.Context, request User) (resData, error) {
	fmt.Println("---------------> Password", request.Password)
	var res resData
	fmt.Println("---------------> Insert")
	err := m.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, userRowsExpectAutoSet)
		fmt.Println(query)

		_, err := s.ExecCtx(ctx, query, request.Name, request.Password, request.Gender)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return err
	})
	if err != nil {
		res.ResCode = "000"
		res.ResCode = "Success User Save."
		return res, err
	}
	return res, err
}
func (m *customUserModel) DeleteUser(ctx context.Context, request int64) (resData, error) {
	fmt.Println(request)
	var res resData
	fmt.Println("=============> Delete")
	data, err := m.FindOne(ctx, request)
	fmt.Println("{{{id}}}", data.Id)
	if err != nil {
		return res, err
	}
	if data != nil {

		userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
		_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
			query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
			fmt.Println("-------------=======>", query)
			return conn.ExecCtx(ctx, query, data.Id)
		}, userIdKey)
	}
	return res, nil
}
