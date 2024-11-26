package wxwork

import (
	"errors"
	"jdy/config"
	"jdy/model"
	"jdy/types"
	"strconv"
	"time"

	"github.com/acmestack/gorm-plus/gplus"
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
			sq, s := gplus.NewQuery[model.Staff]()
			sq.Eq(&s.Phone, detail.Mobile)
			staff, sdb := gplus.SelectOne(sq)
			// 数据库错误
			if sdb.Error != nil && !errors.Is(sdb.Error, gorm.ErrRecordNotFound) {
				return nil, errors.New("注册员工失败")
			}
			// 员工不存在
			if errors.Is(sdb.Error, gorm.ErrRecordNotFound) {
				staff.Phone = &detail.Mobile

				staff.NickName = userinfo.Name
				staff.Avatar = detail.Avatar
				staff.Email = detail.Email

				gender, err := strconv.Atoi(detail.Gender)
				if err != nil {
					gender = 0
				}
				staff.Gender = uint(gender)
				if db := gplus.Insert(staff); db.Error != nil {
					return nil, errors.New("员工注册失败")
				}
			}

			account.StaffId = &staff.Id
			if db := gplus.UpdateById(&account); db.Error != nil {
				return nil, errors.New("更新账号信息失败")
			}
		}
	}

	// 查询员工
	staff, err := logic.getStaff(*account.StaffId)
	if err != nil {
		return nil, err
	}

	// 更新账号
	account.LastLoginIp = ctx.ClientIP()
	now := time.Now()
	account.LastLoginAt = &now

	if db := gplus.UpdateById(&account); db.Error != nil {
		return nil, errors.New("更新账号信息失败")
	}

	return staff, nil
}

// 获取账号
func (wxwork_login_logic) getAccount(uid string) (*model.Account, error) {
	// 查询账号
	aq, a := gplus.NewQuery[model.Account]()
	aq.Eq(&a.Username, uid).And().Eq(&a.Platform, types.PlatformTypeWxWork)
	account, adb := gplus.SelectOne(aq)

	// 账号不存在
	if adb.Error != nil || errors.Is(adb.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("非授权用户禁止登录")
	}

	return account, nil
}

// 获取员工
func (wxwork_login_logic) getStaff(id string) (*model.Staff, error) {
	// 查询账号
	sq, s := gplus.NewQuery[model.Staff]()
	sq.Eq(&s.Id, id)
	staff, sdb := gplus.SelectOne(sq)

	// 账号不存在
	if sdb.Error != nil || errors.Is(sdb.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("员工不存在")
	}

	return staff, nil
}
