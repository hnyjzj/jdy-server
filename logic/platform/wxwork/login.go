package wxwork

import (
	"errors"
	"jdy/config"
	"jdy/model"
	"jdy/types"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type wxwork_login_logic struct{}

// 企业微信授权登录
func (w *WxWorkLogic) Login(ctx *gin.Context, code string) (*model.Staff, error) {
	var (
		jdy = config.NewWechatService().JdyWork

		logic = wxwork_login_logic{}
	)

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取用户信息
	user, err := jdy.OAuth.Provider.GetUserInfo(code)
	if err != nil || user.UserID == "" {
		return nil, errors.New("获取企业微信用户信息失败")
	}
	// 读取企业员工信息
	userinfo, err := jdy.User.Get(ctx, user.UserID)
	if err != nil || userinfo.UserID == "" {
		return nil, errors.New("获取企业微信用户信息失败")
	}

	// 获取账号
	account, err := logic.getAccount(user.UserID)
	if err != nil {
		return nil, err
	}

	// 判断是否是首次登录
	if account.Phone == nil && user.UserTicket == "" {
		return nil, errors.New("首次登录需通过企业微信工作台打开并授权手机号")
	}
	// 判断是否已注册
	if account.StaffId == nil && user.UserTicket == "" {
		return nil, errors.New("首次登录需通过企业微信工作台打开并授权手机号")
	}

	// 获取用户详情
	if user.UserTicket != "" {
		// 获取用户详情
		detail, err := jdy.OAuth.Provider.GetUserDetail(user.UserTicket)
		if err != nil {
			return nil, errors.New("获取企业微信用户详情失败")
		}

		// 获取不到手机号
		if detail.Mobile == "" {
			return nil, errors.New("非企业微信授权用户禁止登录")
		}

		// 账号没有手机号
		if account.Phone == nil {
			account.Phone = &detail.Mobile
		}

		// 手机号不一致
		if *account.Phone != detail.Mobile {
			return nil, errors.New("手机号不一致")
		}

		// 判断是否已注册
		if account.StaffId == nil {
			// 查询员工
			var staff model.Staff
			err := tx.Where(&model.Staff{Phone: &detail.Mobile}).First(&staff).Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("注册员工失败")
				}
			}
			// 员工不存在
			if errors.Is(err, gorm.ErrRecordNotFound) {
				staff.Phone = &detail.Mobile

				staff.Nickname = userinfo.Name
				staff.Avatar = detail.Avatar
				staff.Email = detail.Email

				gender, err := strconv.Atoi(detail.Gender)
				if err != nil {
					gender = 0
				}
				staff.Gender = uint(gender)
				if db := tx.Save(&staff); db.Error != nil {
					tx.Rollback()
					return nil, errors.New("员工注册失败")
				}
			}

			account.StaffId = &staff.Id
			if db := tx.Save(&account); db.Error != nil {
				tx.Rollback()
				return nil, errors.New("注册员工失败，请联系管理员")
			}
		}
	}

	// 查询员工
	staff, err := logic.getStaff(*account.StaffId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新账号
	account.LastLoginIp = ctx.ClientIP()
	now := time.Now()
	account.LastLoginAt = &now

	if db := tx.Save(&account); db.Error != nil {
		tx.Rollback()
		return nil, errors.New("更新账号信息失败")
	}

	return staff, tx.Commit().Error
}

// 获取账号
func (wxwork_login_logic) getAccount(uid string) (*model.Account, error) {
	// 查询账号
	var account model.Account
	if err := model.DB.Where(&model.Account{Username: &uid, Platform: types.PlatformTypeWxWork}).First(&account).Error; err != nil {
		return nil, errors.New("账号不存在")
	}

	return &account, nil
}

// 获取员工
func (wxwork_login_logic) getStaff(id string) (*model.Staff, error) {
	// 查询账号
	var staff model.Staff
	if err := model.DB.Model(&model.Staff{}).First(&staff, id).Error; err != nil {
		return nil, errors.New("员工不存在")
	}

	return &staff, nil
}
