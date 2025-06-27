package message

import (
	"errors"
	"log"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/groupRobot/request"
)

// 消息内容
type CaptureScreenMessage struct {
	Username  string `json:"username"`
	Storename string `json:"storename"`
	Url       string `json:"url"`
}

// 发送截屏事件消息
func (M *BaseMessage) SendCaptureScreenMessage(req *CaptureScreenMessage) error {
	// 消息内容
	if res, err := M.WXWork.GroupRobot.SendTemplateCard(M.Ctx, M.Config.Robot.Warnning, &request.GroupRobotMsgTemplateCard{
		CardType: "text_notice",
		MainTitle: request.TemplateCardMainTitle{
			Title: "截屏警告",
			Desc:  time.Now().Format(time.DateTime),
		},
		EmphasisContent: request.TemplateCardEmphasisContent{
			Title: req.Username,
			Desc:  req.Storename,
		},
		CardAction: request.TemplateCardCardAction{
			Type: 1,
			Url:  req.Url,
		},
	}); err != nil || res.ErrCode != 0 {
		log.Printf("发送消息失败: err=%v, response=%+v\n", err.Error(), res)
		return errors.New("发送消息失败" + res.ErrMsg)
	}

	return nil
}
