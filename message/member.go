package message

import (
	"fmt"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
)

type MemberCreateMessage struct {
	ToUser         string `json:"to_user"`          // 接收消息的用户ID
	ExternalUserID string `json:"external_user_id"` // 外部联系人的ID
	Name           string `json:"name"`             // 外部联系人的名称
}

// 发送新增会员提醒
func (M *BaseMessage) SendMemberCreateMessage(req *MemberCreateMessage) {
	url := fmt.Sprintf("%s/member/lists/edit?external_user_id=%s", M.Config.Jdy.Home, req.ExternalUserID)
	messages := &request.RequestMessageSendTemplateCard{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  req.ToUser, // 接收消息的用户ID列表，列表不可为空，最多支持100个用户
			MsgType: "template_card",
			AgentID: M.Config.Jdy.Id,
		},
		TemplateCard: &request.RequestTemplateCard{
			CardType: "text_notice",
			MainTitle: &request.TemplateCardMainTitle{
				Title: "新增会员提醒",
				Desc:  "会员已创建，请及时更新信息",
			},
			HorizontalContentList: []*request.TemplateCardHorizontalContentListItem{
				{
					Type:    0,
					Keyname: "会员名称",
					Value:   req.Name,
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

	if res, err := M.WXWork.Message.SendTemplateCard(M.Ctx, messages); err != nil || (res != nil && res.ErrCode != 0) {
		log.Println("发送消息失败:", err)
		log.Printf("res: %+v\n", res)
	}
}
