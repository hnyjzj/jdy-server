package user

import (
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"

	"github.com/acmestack/gorm-plus/gplus"
)

type UserLogic struct {
	logic.Base
}

// 创建用户逻辑
func (l *UserLogic) CreateUser(req *types.UserReq) (*types.UserRes, *errors.Errors) {
	// 检查用户名和手机号是否已存在
	query, u := gplus.NewQuery[model.User]()
	query.Eq(&u.Phone, req.Phone).Or().Eq(&u.Username, req.Username)
	hasUser, db := gplus.Exists(query)
	if db.Error != nil {
		return nil, errors.ErrUserNotFound
	}

	// 如果用户已存在，则返回错误
	if hasUser {
		return nil, errors.New("用户已存在")
	}

	// 创建用户
	user := model.User{
		Phone:    &req.Phone,
		Username: &req.Username,
		Password: req.Password,

		NickName: req.NickName,
	}
	result := gplus.Insert(&user)
	if result.Error != nil {
		return nil, errors.New("创建用户失败")
	}

	// 返回用户信息
	res := &types.UserRes{
		Id:       user.Id,
		Phone:    *user.Phone,
		Username: *user.Username,
		NickName: user.NickName,
	}

	return res, nil
}

// 获取用户信息
func (l *UserLogic) GetUserInfo(uid string) (*types.UserRes, error) {
	query, u := gplus.NewQuery[model.User]()
	query.Eq(&u.Id, uid)

	user, db := gplus.SelectGeneric[model.User, *types.UserRes](query)
	if db.Error != nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}
