package wxwork

import (
	"fmt"
	"jdy/config"
	"jdy/enums"
	"jdy/errors"
	"jdy/message"
	"jdy/model"
	"jdy/utils"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type wxworkLoginLogic struct {
	Staff   *model.Staff
	Account *model.Account

	Ctx         *gin.Context
	Db          *gorm.DB
	UserInfo    *model.Account
	Departments []int
}

func (w *WxWorkLogic) CodeLogin(ctx *gin.Context, code string) (*model.Staff, error) {
	l := &wxworkLoginLogic{
		Ctx: ctx,
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Db = tx

		// 获取用户信息
		if err := l.getCodeUserInfo(code); err != nil {
			return err
		}

		// 获取账号
		if err := l.getAccount(false); err != nil {
			return err
		}

		// 判断是否已注册
		if l.Account.StaffId == nil {
			return errors.New("首次登录需通过企业微信工作台打开并授权手机号")
		}

		// 获取员工信息
		if err := l.getStaff(); err != nil {
			return err
		}

		// 更新登录信息
		if err := l.updateLogin(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return l.Staff, nil
}

func (w *WxWorkLogic) OauthLogin(ctx *gin.Context, code string, isRegister bool) (*model.Staff, error) {
	l := &wxworkLoginLogic{
		Ctx: ctx,
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		l.Db = tx

		// 获取用户信息
		if err := l.getOathUserInfo(code); err != nil {
			return err
		}

		// 获取账号
		if err := l.getAccount(isRegister); err != nil {
			return err
		}

		// 判断是否需要注册
		if err := l.register(); err != nil {
			return err
		}

		// 更新账号
		if err := l.updateAccount(); err != nil {
			return err
		}

		// 获取员工信息
		if err := l.getStaff(); err != nil {
			return err
		}

		// 更新登录信息
		if err := l.updateLogin(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return l.Staff, nil
}

func (l *wxworkLoginLogic) getCodeUserInfo(code string) error {

	var (
		jdy = config.NewWechatService().JdyWork
	)

	// 获取用户信息
	user, err := jdy.OAuth.Provider.GetUserInfo(code)
	if err != nil || user.UserID == "" {
		return errors.New("获取企业微信用户信息失败")
	}

	l.UserInfo = &model.Account{
		Platform: enums.PlatformTypeWxWork,
		Username: &user.UserID,
	}

	return nil
}

func (l *wxworkLoginLogic) getOathUserInfo(code string) error {
	var (
		jdy = config.NewWechatService().JdyWork
	)

	// 获取用户信息
	user, err := jdy.OAuth.Provider.GetUserInfo(code)
	if err != nil || user.UserID == "" || user.UserTicket == "" {
		log.Printf("获取企业微信用户信息失败: %+v", err)
		log.Printf("获取企业微信用户信息失败: %+v", user)
		return errors.New("获取企业微信用户信息失败")
	}
	// 获取用户详情
	detail, err := jdy.OAuth.Provider.GetUserDetail(user.UserTicket)
	if err != nil || detail.Mobile == "" {
		log.Printf("获取企业微信用户详情失败: %+v", err)
		log.Printf("获取企业微信用户详情失败: %+v", detail)
		return errors.New("获取企业微信用户详情失败")
	}
	// 读取员工信息
	userinfo, err := jdy.User.Get(l.Ctx, user.UserID)
	if err != nil || userinfo.UserID == "" {
		log.Printf("读取员工信息失败: %+v", err)
		log.Printf("读取员工信息失败: %+v", userinfo)
		return errors.New("读取员工信息失败")
	}

	// 获取性别
	var gender enums.Gender

	l.UserInfo = &model.Account{
		Platform: enums.PlatformTypeWxWork,

		Phone:    &detail.Mobile,
		Username: &user.UserID,
		Nickname: &userinfo.Name,
		Avatar:   &detail.Avatar,
		Email:    &detail.Email,
		Gender:   gender.Convert(detail.Gender),
	}
	l.Departments = userinfo.Department

	return nil
}

// 获取账号
func (l *wxworkLoginLogic) getAccount(isRegister bool) error {
	// 查询账号
	if err := l.Db.Where(&model.Account{
		Platform: enums.PlatformTypeWxWork,
		Username: l.UserInfo.Username,
	}).First(&l.Account).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("查询账号失败")
		}
	}

	// 账号不存在
	if l.Account.Id == "" {
		if !isRegister {
			return errors.New("账号不存在")
		}
		l.Account = &model.Account{
			Platform: enums.PlatformTypeWxWork,
			Username: l.UserInfo.Username,
		}

		if err := l.Db.Save(&l.Account).Error; err != nil {
			return errors.New("账号注册失败")
		}
	}

	// 判断是否是首次登录
	if l.Account.Phone == nil && l.UserInfo.Phone == nil {
		return errors.New("首次登录需通过企业微信工作台打开并授权手机号")
	}

	return nil
}

// 更新账号
func (l *wxworkLoginLogic) updateAccount() error {
	// 更新账号信息
	if err := l.Db.Model(&l.Account).Updates(&model.Account{
		Phone:    l.UserInfo.Phone,
		Nickname: l.UserInfo.Nickname,
		Avatar:   l.UserInfo.Avatar,
		Email:    l.UserInfo.Email,
		Gender:   l.UserInfo.Gender,
	}).Error; err != nil {
		return errors.New("更新账号失败")
	}

	return nil
}

func (l *wxworkLoginLogic) register() error {
	// 账号没有手机号
	if l.Account.Phone == nil {
		l.Account.Phone = l.UserInfo.Phone
	}

	// 手机号未授权
	if l.UserInfo.Phone == nil {
		return errors.New("手机号未授权")
	}

	// 手机号不一致
	if l.Account.Phone != nil && *l.Account.Phone != *l.UserInfo.Phone {
		return errors.New("手机号不一致")
	}

	if l.Account.StaffId != nil {
		return nil
	}

	// 查询员工
	var data *model.Staff
	err := l.Db.Where(&model.Staff{Phone: l.UserInfo.Phone}).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("注册员工失败")
	}

	// 员工不存在
	if data.Id == "" {
		if err := l.createStaff(); err != nil {
			return err
		}
		if err := l.createAccount(); err != nil {
			return err
		}
		if err := l.addRoles(); err != nil {
			return err
		}
		if err := l.addDepartments(); err != nil {
			return err
		}
	}

	l.Account.StaffId = &data.Id

	return nil
}

// 创建员工
func (l *wxworkLoginLogic) createStaff() error {
	// 创建员工
	data := &model.Staff{
		Phone:    l.UserInfo.Phone,
		Nickname: *l.UserInfo.Nickname,
		Avatar:   *l.UserInfo.Avatar,
		Email:    *l.UserInfo.Email,
		Gender:   l.UserInfo.Gender,
	}
	if err := l.Db.Create(&data).Error; err != nil {
		return errors.New("员工注册失败")
	}

	l.Staff = data

	return nil
}

// 创建账号
func (l *wxworkLoginLogic) createAccount() error {
	// 查询账号
	var account *model.Account
	if err := l.Db.Where(&model.Account{
		Platform: enums.PlatformTypeAccount,
		Phone:    l.UserInfo.Phone,
	}).First(&account).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("查询账号失败")
		}
	}

	// 账号不存在
	if account.Id == "" {
		password := utils.RandomAlphanumeric(8)
		account = &model.Account{
			Platform: enums.PlatformTypeAccount,
			Phone:    l.UserInfo.Phone,
			Password: &password,
			Username: l.UserInfo.Username,
			Nickname: l.UserInfo.Nickname,
			Avatar:   l.UserInfo.Avatar,
			Email:    l.UserInfo.Email,
			Gender:   l.UserInfo.Gender,
			StaffId:  &l.Staff.Id,
		}
		if err := l.Db.Create(&account).Error; err != nil {
			return errors.New("创建账号失败")
		}

		go func(UserInfo *model.Account) {
			m := message.NewMessage(l.Ctx)
			m.SendRegisterMessage(&message.RegisterMessageContent{
				Nickname: *UserInfo.Nickname,
				Username: *UserInfo.Username,
				Phone:    *UserInfo.Phone,
				Password: password,
			})
		}(l.UserInfo)
	}

	return nil
}

// 分配权限
func (l *wxworkLoginLogic) addRoles() error {
	// 查询默认权限
	role, err := model.Role{}.Default()
	if err != nil {
		return errors.New("查询默认权限失败")
	}

	if err := l.Db.Model(&l.Staff).Association("Roles").Append(role); err != nil {
		return errors.New("分配权限失败")
	}

	return nil
}

// 分配部门
func (l *wxworkLoginLogic) addDepartments() error {

	var (
		jdy       = config.NewWechatService().JdyWork
		storeIds  []string
		RegionIds []string
	)

	for _, id := range l.Departments {
		// 获取部门信息
		party, err := jdy.Department.Get(l.Ctx, id)
		if err != nil || party.ErrCode != 0 {
			log.Printf("获取部门失败: %+v\n", party)
			if err == nil {
				err = fmt.Errorf("wechat api error: %d %s", party.ErrCode, party.ErrMsg)
			}
			return err
		}

		switch {
		case strings.Contains(party.Department.Name, "店"):
			storeIds = append(storeIds, fmt.Sprint(party.Department.ID))
		case strings.Contains(party.Department.Name, "区域"):
			RegionIds = append(RegionIds, fmt.Sprint(party.Department.ID))
		}
	}

	var stores []model.Store
	if err := l.Db.Where("id_wx in (?)", storeIds).Find(&stores).Error; err != nil {
		return err
	}
	if err := l.Db.Model(&l.Staff).Association("Stores").Append(stores); err != nil {
		return err
	}

	var regions []model.Region
	if err := l.Db.Where("id_wx in (?)", RegionIds).Find(&regions).Error; err != nil {
		return err
	}
	if err := l.Db.Model(&l.Staff).Association("Regions").Append(regions); err != nil {
		return err
	}

	return nil
}

// 获取员工
func (l *wxworkLoginLogic) getStaff() error {
	// 查询账号
	if err := l.Db.
		Preload("Account", func(tx *gorm.DB) *gorm.DB {
			return tx.Where(&model.Account{Platform: enums.PlatformTypeWxWork})
		}).First(&l.Staff, "id = ?", l.Account.StaffId).Error; err != nil {
		return errors.ErrStaffNotFound
	}

	// 判断是否被禁用
	if l.Staff.IsDisabled {
		return errors.ErrStaffDisabled
	}

	return nil
}

// 更新登录信息
func (l *wxworkLoginLogic) updateLogin() error {
	// 更新账号
	l.Account.UpdateLoginInfo(l.Ctx.ClientIP())

	// 更新员工
	if db := l.Db.Save(&l.Account); db.Error != nil {
		return errors.New("登录失败")
	}

	return nil
}
