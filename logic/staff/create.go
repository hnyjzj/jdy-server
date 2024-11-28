package staff

import (
	"fmt"
	"jdy/config"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"strconv"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建员工
func (l *StaffLogic) CreateAccount(ctx *gin.Context, req *types.StaffReq) *errors.Errors {
	// 创建账号
	switch req.Platform {
	case types.PlatformTypeAccount:
		if err := l.account(ctx, req.Account); err != nil {
			return errors.New(err.Error())
		}
	case types.PlatformTypeWxWork:

		if err := l.wxwork(ctx, req.WxWork); err != nil {
			return errors.New(err.Error())
		}
	default:
		return errors.New("平台类型错误")
	}

	return nil
}

// 创建账号（账号密码登录）
func (l *StaffLogic) account(ctx *gin.Context, req *types.StaffAccountReq) error {
	// 启动事务
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询员工
	var staff model.Staff
	if err := tx.
		Unscoped().
		Where(&model.Staff{
			Phone: &req.Phone,
		}).
		First(&staff).
		Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("查询账号失败")
		}
	}

	// 如果账号已存在，则返回错误
	if staff.Id != "" {
		return errors.New("账号已存在")
	}

	// 查询员工
	var account model.Account
	if err := tx.
		Unscoped().
		Where(&model.Account{
			Phone: &req.Phone,
		}).
		Or(&model.Account{
			Username: &req.Username,
		}).
		First(&account).
		Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("查询账号失败")
		}
	}

	// 如果账号已存在，则返回错误
	if account.Id != "" {
		return errors.New("账号已存在")
	}

	// 创建账号
	data := &model.Staff{
		Phone:    &req.Phone,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Email:    req.Email,

		Account: &model.Account{
			Platform: types.PlatformTypeAccount,

			Phone:    &req.Phone,
			Username: &req.Username,
			Password: &req.Password,

			Nickname: &req.Nickname,
			Avatar:   &req.Avatar,
			Email:    &req.Email,
		},
	}

	// 创建账号
	if err := tx.Save(data).Error; err != nil {
		tx.Rollback()
		return errors.New("创建账号失败")
	}

	go func() {
		l.sendCreateMessage(ctx, &MessageContent{
			Nickname: data.Nickname,
			Username: req.Username,
			Phone:    req.Phone,
			Password: req.Password,
		})
	}()

	return tx.Commit().Error
}

func (l *StaffLogic) wxwork(ctx *gin.Context, req *types.StaffWxWorkReq) error {
	var (
		wxwork = config.NewWechatService().JdyWork
	)

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 循环 req.userid 获取企业微信用户信息
	for _, userid := range req.UserId {
		// 获取企业微信用户信息
		user, err := wxwork.User.Get(ctx, fmt.Sprint(userid))
		if err != nil {
			return errors.New(fmt.Sprintf("获取企业微信用户信息失败：%s", userid))
		}

		// 根据用户名检查账号是否已存在
		var account model.Account
		if err := tx.
			Unscoped().
			Where(&model.Account{
				Username: &user.UserID,
				Platform: types.PlatformTypeWxWork,
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

		gender, err := strconv.Atoi(user.Gender)
		if err != nil {
			gender = 0
		}
		acc.Gender = uint(gender)

		if err := tx.Save(acc).Error; err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("创建账号失败: %s", userid))
		}

		go func() {
			l.sendCreateMessage(ctx, &MessageContent{
				Nickname: user.Name,
				Username: user.UserID,
				Phone:    "暂未授权",
				Password: "无",
			})
		}()
	}

	return tx.Commit().Error
}

// 消息内容
type MessageContent struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (l *StaffLogic) sendCreateMessage(ctx *gin.Context, message *MessageContent) error {
	var (
		wxwork = config.NewWechatService().JdyWork
	)

	var configTemplate string = strings.Join([]string{
		"欢迎加入金斗云 ！",
		">昵  称：%v",
		">账  号：%v",
		">手机号：%v",
		">密  码：%v",
		"",
		">请妥善保管账号信息。",
		">如果手机号未授权，请通过其他方式登录。",
		">如有疑问，请联系管理员。",
	}, "\n")
	content := fmt.Sprintf(configTemplate,
		message.Nickname,
		message.Username,
		message.Phone,
		message.Password,
	)

	messages := &request.RequestMessageSendMarkdown{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  message.Username,
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
