package message

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
)

type CustomerStatisticsPersonal struct {
	ToUser string `json:"to_user"` // 接收消息的用户

	StatTime            time.Time `json:"stat_time"`             // 数据日期
	ChatCnt             int64     `json:"chat_cnt"`              // 聊天总数
	MessageCnt          int64     `json:"message_cnt"`           // 发送消息数
	ReplyPercentage     float64   `json:"reply_percentage"`      // 已回复聊天占比
	AvgReplyTime        int64     `json:"avg_reply_time"`        // 平均首次回复时长(分钟)
	NegativeFeedbackCnt int64     `json:"negative_feedback_cnt"` // 删除/拉黑成员的客户数
	NewApplyCnt         int64     `json:"new_apply_cnt"`         // 发起申请数
	NewContactCnt       int64     `json:"new_contact_cnt"`       // 新增客户数
}

// 发送新增会员提醒
func (M *BaseMessage) SendCustomerStatisticsPersonal(req *CustomerStatisticsPersonal) {
	messages := &request.RequestMessageSendMarkdown{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  req.ToUser,
			MsgType: "markdown",
			AgentID: M.Config.Jdy.Id,
		},
		Markdown: &request.RequestMarkdown{
			Content: strings.Join([]string{
				"**客户统计**",
				"**数据日期**：" + req.StatTime.Format("2006-01-02"),
				"**发起申请数**：" + fmt.Sprint(req.NewApplyCnt),
				"**新增客户数**：" + fmt.Sprint(req.NewContactCnt),
				"**聊天总数**：" + fmt.Sprint(req.ChatCnt),
				"**发送消息数**：" + fmt.Sprint(req.MessageCnt),
				"**已回复聊天占比**：" + fmt.Sprintf("%.2f", req.ReplyPercentage) + "%",
				"**平均首次回复时长**：" + fmt.Sprint(req.AvgReplyTime) + "分钟",
				"**删除/拉黑成员的客户数**：" + fmt.Sprint(req.NegativeFeedbackCnt),
			}, "\n"),
		},
	}

	if res, err := M.WXWork.Message.SendMarkdown(M.Ctx, messages); err != nil || (res != nil && res.ErrCode != 0) {
		log.Println("发送消息失败:", err)
		log.Printf("res: %+v\n", res)
	}
}
