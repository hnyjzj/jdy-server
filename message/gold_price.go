package message

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
)

type GoldPriceMessage struct {
	ToUser    []string `json:"to_user"`    // 接收消息的用户ID列表，列表不可为空，最多支持100个用户
	StoreName string   `json:"store_name"` // 门店名称
	Operator  string   `json:"operator"`   // 操作人
}

// 发送黄金价格设置提醒
func (M *BaseMessage) SendGoldPriceSetMessage(req *GoldPriceMessage) {
	url := fmt.Sprintf("%s/system/gold/price", M.App.Home)
	ToUser := strings.Join(req.ToUser, "|")
	messages := &request.RequestMessageSendTemplateCard{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  ToUser, // 接收消息的用户ID列表，列表不可为空，最多支持100个用户
			MsgType: "template_card",
			AgentID: M.App.Id,
		},
		TemplateCard: &request.RequestTemplateCard{
			CardType: "text_notice",
			MainTitle: &request.TemplateCardMainTitle{
				Title: "【今日金价】设置提醒",
				Desc:  "请及时更新今日金价",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "门店名称",
					Value:   req.StoreName,
				},
			},
			CardAction: &request.TemplateCardAction{
				Type: 1,
				Url:  url,
			},
			JumpList: []*request.TemplateCardJumpListItem{
				{
					Type:  1,
					Title: "查看详情",
					Url:   url,
				},
			},
		},
	}

	if a, err := M.WXWork.Message.SendTemplateCard(M.Ctx, messages); err != nil {
		log.Println("发送消息失败:", err)
		fmt.Printf("a: %+v\n", a)
	}
}

// 发送黄金价格更新提醒
func (M *BaseMessage) SendGoldPriceUpdateMessage(req *GoldPriceMessage) {
	url := fmt.Sprintf("%s/system/gold/price", M.App.Home)
	ToUser := strings.Join(req.ToUser, "|")
	messages := &request.RequestMessageSendTemplateCard{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  ToUser,
			MsgType: "template_card",
			AgentID: M.App.Id,
		},
		TemplateCard: &request.RequestTemplateCard{
			CardType: "text_notice",
			MainTitle: &request.TemplateCardMainTitle{
				Title: "【今日金价】变更通知",
				Desc:  "黄金价格已变更，请及时查收",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "门店名称",
					Value:   req.StoreName,
				},
				{
					Type:    0,
					Keyname: "操作人",
					Value:   req.Operator,
				},
				{
					Type:    0,
					Keyname: "变更时间",
					Value:   time.Now().Format("2006-01-02 15:04"),
				},
			},
			CardAction: &request.TemplateCardAction{
				Type: 1,
				Url:  url,
			},
			JumpList: []*request.TemplateCardJumpListItem{
				{
					Type:  1,
					Title: "查看详情",
					Url:   url,
				},
			},
		},
	}

	if a, err := M.WXWork.Message.SendTemplateCard(M.Ctx, messages); err != nil {
		log.Println("发送消息失败:", err)
		fmt.Printf("a: %+v\n", a)
	}
}
