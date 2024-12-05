package wxwork

import (
	"fmt"
	"jdy/config"
	"log"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/gin-gonic/gin"
)

// 消息内容
type RegisterMessageContent struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func SendRegisterMessage(ctx *gin.Context, message *RegisterMessageContent) {
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
	if err != nil {
		log.Println("发送消息失败:", err)
	}
}
