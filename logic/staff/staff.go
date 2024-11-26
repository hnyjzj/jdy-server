package staff

import (
	"encoding/json"
	"fmt"
	"jdy/config"
	"jdy/errors"
	"jdy/logic"
	"jdy/model"
	"jdy/types"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StaffLogic struct {
	logic.Base
}

// 创建员工
func (l *StaffLogic) CreateAccount(ctx *gin.Context, req *types.StaffReq) *errors.Errors {

	var (
		err *errors.Errors
	)

	// 创建账号
	switch req.Platform {
	case types.PlatformTypeAccount:
		err = l.createAccount(ctx, req.Account)
		if err != nil {
			return err
		}
	case types.PlatformTypeWxWork:
		err = l.createWxWork(ctx, req.WxWork)
		if err != nil {
			return err
		}
	default:
		return errors.New("平台类型错误")
	}

	return nil
}

// 创建账号（账号密码登录）
func (l *StaffLogic) createAccount(ctx *gin.Context, req *types.StaffAccountReq) *errors.Errors {
	// 查询账号
	aq, a := gplus.NewQuery[model.Account]()
	aq.Eq(&a.Platform, types.PlatformTypeAccount).And().
		Eq(&a.Phone, req.Phone)
	account, resultDb := gplus.SelectOne(aq)
	if resultDb.Error != nil && !errors.Is(resultDb.Error, gorm.ErrRecordNotFound) {
		return errors.New("查询账号失败")
	}
	// 如果账号已存在，则返回错误
	if account != nil {
		return errors.New("账号已存在")
	}

	// 创建账号
	account = &model.Account{
		Platform: types.PlatformTypeAccount,

		Phone:    &req.Phone,
		Username: &req.Username,
		Password: &req.Password,
	}

	// 根据手机号检查员工是否已存在
	sq, s := gplus.NewQuery[model.Staff]()
	sq.Eq(&s.Phone, req.Phone)
	staff, resultDb := gplus.SelectOne(sq)
	if resultDb.Error != nil && !errors.Is(resultDb.Error, gorm.ErrRecordNotFound) {
		return errors.New("查询账号失败")
	}
	// 如果员工不存在，则创建员工
	if staff == nil {
		staff = &model.Staff{
			Phone:    &req.Phone,
			NickName: req.NickName,
		}
		result := gplus.Insert(&staff)
		if result.Error != nil {
			return errors.New("创建员工失败")
		}

		account.StaffId = &staff.Id
	}

	// 创建账号
	result := gplus.Insert(&account)
	if result.Error != nil {
		return errors.New("创建账号失败")
	}

	go l.sendCreateMessage(ctx, account, staff)

	return nil
}

func (l *StaffLogic) createWxWork(ctx *gin.Context, req *types.StaffWxWorkReq) *errors.Errors {
	var (
		wxwork = config.NewWechatService().JdyWork
	)
	// 循环 req.userid 获取企业微信用户信息
	for _, userid := range req.UserId {
		// 获取企业微信用户信息
		user, err := wxwork.User.Get(ctx, fmt.Sprint(userid))
		if err != nil {
			return errors.New(fmt.Sprintf("获取企业微信用户信息失败：%s", userid))
		}

		// 根据用户名检查账号是否已存在
		aq, a := gplus.NewQuery[model.Account]()
		aq.Eq(&a.Username, user.UserID).And().Eq(&a.Platform, types.PlatformTypeWxWork)
		account, resultDb := gplus.SelectOne(aq)
		if resultDb.Error != nil && !errors.Is(resultDb.Error, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("查询账号失败: %s", userid))
		}
		// 如果账号已存在，则跳过循环
		if account.Id != "" {
			continue
		}
		// 创建账号
		account = &model.Account{
			Platform: types.PlatformTypeWxWork,
			Username: &user.UserID,

			Nickname: &user.Name,
			Avatar:   &user.Avatar,
			Email:    &user.Email,
		}

		if result := gplus.Insert(&account); result.Error != nil {
			return errors.New(fmt.Sprintf("创建账号失败: %s", userid))
		}

		phone := "暂未授权"
		passwprd := "无"

		staff := &model.Staff{
			Phone:    &phone,
			NickName: user.Name,
		}
		account.Password = &passwprd

		go l.sendCreateMessage(ctx, account, staff)
	}

	return nil
}

func (l *StaffLogic) sendCreateMessage(ctx *gin.Context, account *model.Account, staff *model.Staff) error {
	var (
		wxwork = config.NewWechatService().JdyWork
	)

	var configTemplate string = strings.Join([]string{
		"欢迎加入金斗云 ！",
		">昵  称：%s",
		">账  号：%s",
		">手机号：%s",
		">密  码：%s",
		"",
		">请妥善保管账号信息。",
		">如果手机号未授权，请通过其他方式登录。",
		">如有疑问，请联系管理员。",
	}, "\n")
	content := fmt.Sprintf(configTemplate,
		staff.NickName,
		*account.Username,
		*staff.Phone,
		*account.Password,
	)

	messages := &request.RequestMessageSendMarkdown{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  *account.Username,
			MsgType: "markdown",
			AgentID: config.Config.Wechat.Work.Jdy.Id,
		},
		Markdown: &request.RequestMarkdown{
			Content: content,
		},
	}
	_, err := wxwork.Message.SendMarkdown(ctx, messages)

	return err
}

// 获取员工信息
func (l *StaffLogic) GetStaffInfo(ctx *gin.Context, user *model.Staff) (*types.StaffRes, error) {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return nil, errors.New("获取员工信息失败")
	}

	var saff types.StaffRes
	if err := json.Unmarshal(userBytes, &saff); err != nil {
		return nil, errors.New("获取员工信息失败")
	}

	return &saff, nil
}
