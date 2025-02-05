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
func (StaffLogic) StaffCreate(ctx *gin.Context, req *types.StaffReq) *errors.Errors {
	l := &AccountCreateLogic{
		Ctx: ctx,
		Req: req,
		Db:  model.DB.Begin(),
	}
	defer func() {
		if r := recover(); r != nil {
			l.Db.Rollback()
		}
	}()

	// 创建账号
	switch l.Req.Platform {
	case types.PlatformTypeAccount:
		if err := l.account(); err != nil {
			l.Db.Rollback()
			return errors.New(err.Error())
		}
	case types.PlatformTypeWxWork:

		if err := l.wxwork(); err != nil {
			l.Db.Rollback()
			return errors.New(err.Error())
		}
	default:
		return errors.New("平台类型错误")
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
		tx  = l.Db
		ctx = l.Ctx
	)

	// 查询账号存不存在
	var account model.Account
	if err := tx.Unscoped().
		Where(&model.Account{
			Platform: types.PlatformTypeAccount,
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
	if err := tx.Unscoped().
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
		Platform: types.PlatformTypeAccount,

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
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
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

	return tx.Commit().Error
}

func (l *AccountCreateLogic) wxwork() error {
	var (
		ctx = l.Ctx
		req = l.Req.WxWork
		tx  = l.Db

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
		if err := tx.
			Unscoped().
			Where(&model.Account{
				Platform: types.PlatformTypeWxWork,
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
			Platform: types.PlatformTypeWxWork,
			Username: &user.UserID,

			Nickname: &user.Name,
			Avatar:   &user.Avatar,
			Email:    &user.Email,
		}

		var gender enums.Gender
		acc.Gender = gender.Convert(user.Gender)

		if err := tx.Save(acc).Error; err != nil {
			tx.Rollback()
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

	return tx.Commit().Error
}
