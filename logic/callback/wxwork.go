package callback

import (
	"errors"
	"jdy/config"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/gin-gonic/gin"
)

type WxWork struct {
	Ctx    *gin.Context            // 上下文
	Event  contract.EventInterface // 事件
	Staff  *types.Staff            // 员工信息
	Wechat *config.WechatService   // 金斗云
}

func NewWxWork(ctx *gin.Context, event contract.EventInterface) *WxWork {
	return &WxWork{
		Ctx:    ctx,
		Event:  event,
		Wechat: config.NewWechatService(),
	}
}

func (w *WxWork) GetStaff() error {
	username := w.Event.GetFromUserName()

	// 查询用户信息
	var account model.Account
	if err := model.DB.Where(&model.Account{
		Username: &username,
		Platform: enums.PlatformTypeWxWork,
	}).First(&account).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 判断用户
	if account.StaffId == nil {
		return errors.New("员工不存在")
	}

	w.Staff = &types.Staff{
		Id: *account.StaffId,
	}

	return nil
}
