package message

import (
	"fmt"
	"jdy/enums"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/shopspring/decimal"
)

type GoldPriceApprovalMessage struct {
	Id        string          // 任务ID
	Price     decimal.Decimal // 价格
	Initiator string          // 发起人
}

// 发送黄金价格审批
func (M *BaseMessage) SendGoldPriceApprovalMessage(req *GoldPriceApprovalMessage) {
	messages := &request.RequestMessageSendTemplateCard{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  "@all",
			MsgType: "template_card",
			AgentID: M.App.Id,
		},
		TemplateCard: &request.RequestTemplateCard{
			CardType: "button_interaction",
			MainTitle: &request.TemplateCardMainTitle{
				Title: "黄金价格变更申请",
				Desc:  "黄金价格有修改申请，请及时审批",
			},
			SubTitleText: fmt.Sprintf("%s 元/克", req.Price.Round(2).String()),
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "发起人",
					Value:   req.Initiator,
				},
			},
			TaskID: req.Id,
			ButtonList: []*request.TemplateCardButtonListItem{
				{
					Text:  "驳回",
					Key:   string(enums.GoldPriceReviewRejected),
					Style: 3,
				},
				{
					Text:  "审批",
					Key:   string(enums.GoldPriceReviewApproved),
					Style: 1,
				},
			},
		},
	}

	if a, err := M.WXWork.Message.SendTemplateCard(M.Ctx, messages); err != nil {
		log.Println("发送消息失败:", err)
		fmt.Printf("a: %+v\n", a)
	}
}

type GoldPriceMessage struct {
	Price     decimal.Decimal // 价格
	Initiator string          // 发起人
	Approver  string          // 审批人
}

// 发送黄金价格消息
func (M *BaseMessage) SendGoldPriceMessage(req *GoldPriceMessage) {
	messages := &request.RequestMessageSendTemplateCard{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  "@all",
			MsgType: "template_card",
			AgentID: M.App.Id,
		},
		TemplateCard: &request.RequestTemplateCard{
			CardType: "text_notice",
			MainTitle: &request.TemplateCardMainTitle{
				Title: "黄金价格变更通知",
				Desc:  "黄金价格已变更，请及时查收",
			},
			EmphasisContent: &request.TemplateCardEmphasisContent{
				Title: req.Price.Round(2).String(),
				Desc:  "元/克",
			},
			CardAction: &request.TemplateCardAction{
				Type: 1,
				Url:  M.App.Home,
			},
			// SubTitleText: "请及时查收",
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "发起人",
					Value:   req.Initiator,
				},
				{
					Type:    0,
					Keyname: "审批人",
					Value:   req.Approver,
				},
			},
			JumpList: []*request.TemplateCardJumpListItem{
				{
					Type:  1,
					Title: "查看详情",
					Url:   M.App.Home,
				},
			},
		},
	}

	if a, err := M.WXWork.Message.SendTemplateCard(M.Ctx, messages); err != nil {
		log.Println("发送消息失败:", err)
		fmt.Printf("a: %+v\n", a)
	}
}

type UpdateGoldPriceMessage struct {
	Code    string // 任务code
	Message string // 任务状态
}

// 更新黄金审批消息
func (M *BaseMessage) UpdateGoldPriceMessage(req *UpdateGoldPriceMessage) {
	messages := power.HashMap{
		"atall":         1,
		"agentid":       M.App.Id,
		"response_code": req.Code,
		"button": power.HashMap{
			"replace_name": req.Message,
		},
	}

	if a, err := M.WXWork.Message.UpdateTemplateCard(M.Ctx, &messages); err != nil {
		log.Println("发送消息失败:", err)
		fmt.Printf("a: %+v\n", a)
	}
}
