package wxwork

import (
	"errors"
	"jdy/config"
	"jdy/model"
	"log"
)

// 获取授权链接
func (w *WxWorkLogic) GetUser(user_id string) (*model.Staff, error) {
	wxwork := config.NewWechatService().JdyWork

	// 获取企业微信用户信息
	user, err := wxwork.User.Get(w.Ctx, user_id)
	if err != nil || user.UserID == "" {
		log.Printf("读取员工信息失败: %+v, %+v", err, user)
		return nil, errors.New("读取员工信息失败")
	}

	var staff model.Staff
	if err := model.DB.Where(&model.Staff{Username: user.UserID}).Attrs(model.Staff{
		Username: user.UserID,
		Nickname: user.Name,
	}).FirstOrInit(&staff).Error; err != nil {
		return nil, err
	}

	return &staff, nil
}
