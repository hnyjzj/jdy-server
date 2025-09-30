package message

import (
	"errors"
	"log"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/groupRobot/request"
)

// 消息内容
type CaptureScreenMessage struct {
	Type string `json:"type"`

	Username  string `json:"username"`
	Storename string `json:"storename"`
	Url       string `json:"url"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
}

// 发送截屏事件消息
func (M *BaseMessage) SendCaptureScreenMessage(req *CaptureScreenMessage) error {
	if M.Config.Robot.Warning == "" {
		return errors.New("请先配置机器人")
	}

	data := &request.GroupRobotMsgTemplateCard{
		CardType: "text_notice",
		MainTitle: request.TemplateCardMainTitle{
			Title: req.Type,
		},
		EmphasisContent: request.TemplateCardEmphasisContent{
			Title: req.Title,
			Desc:  time.Now().Format(time.DateTime),
		},
		CardAction: request.TemplateCardCardAction{
			Type: 1,
			Url:  req.Url,
		},
		HorizontalContentList: []request.TemplateCardHorizontalContentListItem{
			{
				KeyName: "姓名",
				Value:   req.Username,
			},
			{
				KeyName: "门店",
				Value:   req.Storename,
			},
		},
	}

	if req.Desc != "" {
		data.SubTitleText = req.Desc
	}

	// 消息内容
	if res, err := M.WXWork.GroupRobot.SendTemplateCard(M.Ctx, M.Config.Robot.Warning, data); err != nil || (res != nil && res.ErrCode != 0) {
		log.Printf("发送消息失败: err=%v, response=%+v\n", err.Error(), res)
		return errors.New("发送消息失败" + res.ErrMsg)
	}

	return nil
}
