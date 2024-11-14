package userlogic

import (
	"jdy/logic"
	"jdy/logic_error"
	usermodel "jdy/model/user"
	usertype "jdy/types/user"

	"github.com/acmestack/gorm-plus/gplus"
)

type UserLogic struct {
	logic.Base
}

// 获取用户信息
func (l *UserLogic) GetUserInfo(uid string) (*usertype.UserRes, error) {
	query, u := gplus.NewQuery[usermodel.User]()
	query.Eq(&u.Id, uid)

	user, db := gplus.SelectGeneric[usermodel.User, *usertype.UserRes](query)
	if db.Error != nil {
		return nil, logic_error.ErrUserNotFound
	}

	return user, nil
}
