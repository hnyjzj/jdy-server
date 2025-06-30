package staff

import (
	"jdy/errors"
	"jdy/logic/platform/wxwork"
	"jdy/message"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建员工
func (StaffLogic) StaffCreate(ctx *gin.Context, req *types.StaffReq) error {
	l := &StaffCreateLogic{
		Ctx: ctx,
		Req: req,
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Db = tx

		// 查询员工是否存在
		if err := l.getStaff(); err != nil {
			return err
		}

		// 查询企业微信
		if err := l.getWechat(); err != nil {
			return err
		}

		// 创建账号
		if err := l.createStaff(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

type StaffCreateLogic struct {
	Ctx *gin.Context
	Req *types.StaffReq
	Db  *gorm.DB

	Staff *model.Staff
}

// 查询员工是否存在
func (l *StaffCreateLogic) getStaff() error {
	var (
		db    = l.Db
		staff model.Staff
	)

	// 根据手机号或用户名查询账号
	db = db.Unscoped()
	db = db.Where(&model.Staff{
		Phone: l.Req.Phone,
	})
	db = db.Or(&model.Staff{
		Username: l.Req.Username,
	})
	if err := db.First(&staff).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("查询账号失败")
		}
	}

	if staff.Id != "" {
		return errors.New("账号已存在")
	}

	return nil
}

// 查询企业微信
func (l *StaffCreateLogic) getWechat() error {
	wxlogic := &wxwork.WxWorkLogic{
		Ctx: l.Ctx,
	}
	user, err := wxlogic.GetUser(l.Req.Username)
	if err != nil || user.Username == "" {
		return errors.New("企业微信用户不存在")
	}

	return nil
}

// 创建账号
func (l *StaffCreateLogic) createStaff() error {
	// 创建账号
	l.Staff = &model.Staff{
		Username: l.Req.Username,
		Phone:    l.Req.Phone,
		Password: l.Req.Password,

		Nickname: l.Req.Nickname,
		Avatar:   l.Req.Avatar,
		Email:    l.Req.Email,
		Gender:   l.Req.Gender,
	}

	if err := l.Db.Create(&l.Staff).Error; err != nil {
		return errors.New("创建账号失败")
	}

	return nil
}

// 发送消息
func (l *StaffCreateLogic) sendMessage() error {

	go func() {
		m := message.NewMessage(l.Ctx)
		m.SendRegisterMessage(&message.RegisterMessageContent{
			Username: l.Req.Username,
			Nickname: l.Req.Nickname,
			Phone:    l.Req.Phone,
			Password: l.Req.Password,
		})
	}()

	return nil
}
