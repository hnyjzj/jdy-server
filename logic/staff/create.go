package staff

import (
	"fmt"
	"jdy/config"
	"jdy/enums"
	"jdy/errors"
	"jdy/message"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建员工
func (StaffLogic) StaffCreate(ctx *gin.Context, req *types.StaffReq) error {
	l := &AccountCreateLogic{
		Ctx: ctx,
		Req: req,
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Db = tx

		// 创建账号
		switch l.Req.Platform {
		case enums.PlatformTypeAccount:
			if err := l.account(); err != nil {
				return err
			}
		case enums.PlatformTypeWxWork:

			if err := l.wxwork(); err != nil {
				return err
			}
		default:
			return errors.New("平台类型错误")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

type AccountCreateLogic struct {
	Ctx *gin.Context
	Req *types.StaffReq
	Db  *gorm.DB
}

// 创建账号（账号密码登录）
func (l *AccountCreateLogic) account() error {
	var (
		req = l.Req.Account
		ctx = l.Ctx
	)

	// 查询账号存不存在
	var account model.Account
	if err := l.Db.Unscoped().
		Where(&model.Account{
			Platform: enums.PlatformTypeAccount,
			Phone:    &req.Phone,
		}).
		First(&account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("查询账号失败")
		}
	}
	// 如果账号已存在，则返回错误
	if account.Id != "" {
		return errors.New("账号已存在")
	}

	// 查询手机号是否已注册
	if err := l.Db.Unscoped().
		Where(&model.Staff{
			Phone: &req.Phone,
		}).
		First(&account.Staff).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("查询账号失败")
		}
	}
	// 如果手机号已注册，则返回错误
	if account.Staff.Id != "" {
		return errors.New("手机号已注册")
	}

	// 创建账号
	account = model.Account{
		Platform: enums.PlatformTypeAccount,

		Phone:    &req.Phone,
		Password: &req.Password,

		Nickname: &req.Nickname,
		Avatar:   &req.Avatar,
		Email:    &req.Email,
		Gender:   req.Gender,

		Staff: &model.Staff{
			Phone:    &req.Phone,
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
			Email:    req.Email,
			Gender:   req.Gender,
		},
	}

	// 加密密码
	if req.Password != "" {
		password, err := account.HashPassword(&req.Password)
		if err != nil {
			return err
		}
		account.Password = &password
	}

	// 创建账号
	if err := l.Db.Save(&account).Error; err != nil {
		return errors.New("创建账号失败")
	}

	go func() {
		m := message.NewMessage(ctx)
		m.SendRegisterMessage(&message.RegisterMessageContent{
			Nickname: req.Nickname,
			Phone:    req.Phone,
			Password: req.Password,
		})
	}()

	return nil
}

func (l *AccountCreateLogic) wxwork() error {
	var (
		ctx = l.Ctx
		req = l.Req.WxWork

		jdy = config.NewWechatService().JdyWork
	)

	// 循环 req.userid 获取企业微信用户信息
	for _, userid := range req.UserId {
		// 获取企业微信用户信息
		user, err := jdy.User.Get(ctx, fmt.Sprint(userid))
		if err != nil || user.UserID == "" {
			return errors.New(fmt.Sprintf("获取企业微信用户信息失败：%s", userid))
		}

		// 根据用户名检查账号是否已存在
		var account model.Account
		if err := l.Db.Unscoped().
			Where(&model.Account{
				Platform: enums.PlatformTypeWxWork,
				Username: &user.UserID,
			}).
			First(&account).
			Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("查询账号失败: %s", userid))
		}
		// 如果账号已存在，则跳过循环
		if account.Id != "" {
			continue
		}

		// 创建账号
		acc := &model.Account{
			Platform: enums.PlatformTypeWxWork,
			Username: &user.UserID,

			Nickname: &user.Name,
			Avatar:   &user.Avatar,
			Email:    &user.Email,
		}

		var gender enums.Gender
		acc.Gender = gender.Convert(user.Gender)

		if err := l.Db.Save(acc).Error; err != nil {
			return errors.New(fmt.Sprintf("创建账号失败: %s", userid))
		}

		go func() {
			m := message.NewMessage(ctx)
			m.SendRegisterMessage(&message.RegisterMessageContent{
				Nickname: user.Name,
				Username: user.UserID,
				Phone:    "暂未授权",
				Password: "无",
			})
		}()
	}

	return nil
}
