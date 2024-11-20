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
