package wxwork

import (
	"errors"
	"jdy/config"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/user/response"
)

// 获取授权链接
func (w *WxWorkLogic) GetUser(user_id string) (*response.ResponseGetUserDetail, error) {
	wxwork := config.NewWechatService().JdyWork

	// 获取企业微信用户信息
	user, err := wxwork.User.Get(w.Ctx, user_id)
	if err != nil || user.UserID == "" {
		log.Printf("读取员工信息失败: %+v, %+v", err, user)
		return nil, errors.New("读取员工信息失败")
	}

	return user, nil
}
