package message

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/shopspring/decimal"
)

// 消息内容
type ReportStatisticMessage struct {
	ToUser []string `json:"to_user"` // 接收者

	StoreName       string    `json:"store_name"`       // 门店名称
	StatisticalTime time.Time `json:"statistical_time"` // 统计时间

	TodayFinished    decimal.Decimal `json:"today_finished"`    // 成品业绩
	TodayOld         decimal.Decimal `json:"today_old"`         // 旧料业绩
	TodayAcciessorie decimal.Decimal `json:"today_acciessorie"` // 配件业绩

	TodayFinisheds map[string]decimal.Decimal `json:"today_finisheds"` // 成品统计

	MonthFinished    decimal.Decimal `json:"month_finished"`    // 月度成品业绩
	MonthOld         decimal.Decimal `json:"month_old"`         // 月度旧料业绩
	MonthAcciessorie decimal.Decimal `json:"month_acciessorie"` // 月度配件业绩
}

// 发送统计报告消息
func (M *BaseMessage) SendReportStatisticMessage(req *ReportStatisticMessage) {

	content := []string{
		fmt.Sprintf("###### %s(%s)", req.StoreName, req.StatisticalTime.Format("01-02")),
		"",
		fmt.Sprintf("销售业绩：**%s**", req.TodayFinished.StringFixed(2)),
		fmt.Sprintf("旧料抵扣：**%s**", req.TodayOld.StringFixed(2)),
		fmt.Sprintf("配件收款：**%s**", req.TodayAcciessorie.StringFixed(2)),
		"",
	}
	for class, total := range req.TodayFinisheds {
		content = append(content,
			fmt.Sprintf("> %s ：**%s**", class, total.StringFixed(2)),
		)
	}
	content = append(content, []string{
		"",
		fmt.Sprintf("月度销售：**%s**", req.MonthFinished.StringFixed(2)),
		fmt.Sprintf("月度抵扣：**%s**", req.MonthOld.StringFixed(2)),
		fmt.Sprintf("月度配件：**%s**", req.MonthAcciessorie.StringFixed(2)),
	}...)

	messages := &request.RequestMessageSendMarkdown{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:  strings.Join(req.ToUser, "|"),
			MsgType: "markdown",
			AgentID: M.Config.Jdy.Id,
		},
		Markdown: &request.RequestMarkdown{
			Content: strings.Join(content, "\n"),
		},
	}

	if res, err := M.WXWork.Message.SendMarkdown(M.Ctx, messages); err != nil || (res != nil && res.ErrCode != 0) {
		log.Println("发送消息失败:", err)
		log.Printf("res: %+v\n", res)
	}
}
