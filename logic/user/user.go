package user

import (
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
)

type UserLogic struct {
	logic.Base
}

// 创建用户逻辑
func (l *UserLogic) CreateUser(ctx *gin.Context, req *types.UserReq) (*types.UserRes, *errors.Errors) {
	var (
		wxwork = config.NewWechatService().JdyWork
	)

	// 检查用户名和手机号是否已存在
	query, u := gplus.NewQuery[model.User]()
	query.Eq(&u.Phone, req.Phone).Or().Eq(&u.Username, req.Username)
	hasUser, db := gplus.Exists(query)
	if db.Error != nil {
		return nil, errors.ErrUserNotFound
	}

	// 如果用户已存在，则返回错误
	if hasUser {
		return nil, errors.New("用户已存在")
	}

	// 创建用户
	user := model.User{
		Phone:    &req.Phone,
		Username: &req.Username,
		Password: req.Password,

		NickName: req.NickName,
	}
	result := gplus.Insert(&user)
	if result.Error != nil {
		return nil, errors.New("创建用户失败")
	}

	// 返回用户信息
	res := &types.UserRes{
		Id:       user.Id,
		Phone:    *user.Phone,
		Username: *user.Username,
		NickName: user.NickName,
	}

	func() {
		var configTemplate string = strings.Join([]string{
			"欢迎加入金斗云，您的账号已创建成功！",
			">昵  称：%s",
			">用户名：%s",
			">手机号：%s",
			">密  码：%s",
			"",
			">请及时修改默认密码，并妥善保管账号信息。",
		}, "\n")
		content := fmt.Sprintf(configTemplate,
			user.NickName,
			*user.Username,
			*user.Phone,
			user.Password,
		)
		messages := &request.RequestMessageSendMarkdown{
			RequestMessageSend: request.RequestMessageSend{
				ToUser:  *user.Username,
				MsgType: "markdown",
				AgentID: config.Config.Wechat.Work.Jdy.Id,
			},
			Markdown: &request.RequestMarkdown{
				Content: content,
			},
		}
		wxwork.Message.SendMarkdown(ctx, messages)
	}()

	return res, nil
}

// 获取用户信息
func (l *UserLogic) GetUserInfo(uid string) (*types.UserRes, error) {
	query, u := gplus.NewQuery[model.User]()
	query.Eq(&u.Id, uid)

	user, db := gplus.SelectGeneric[model.User, *types.UserRes](query)
	if db.Error != nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}
