package setting

import (
	"jdy/config"
	"log"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/agent/request"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type WorkbenchTemplate struct {
	Price decimal.Decimal `json:"price"` // 今日黄金价格
}

// 获取工作台模板
func SetWorkbenchTemplate(ctx *gin.Context, template WorkbenchTemplate) {
	var (
		wxwork = config.NewWechatService().JdyWork
		app    = &config.Config.Wechat.Work.Jdy
	)

	options := &request.RequestSetWorkbenchTemplate{
		AgentID:         app.Id,
		Type:            "keydata",
		ReplaceUserData: true,
		KeyData: request.WorkBenchKeyData{
			Items: []request.WorkBenchKeyDataItem{{
				Key:  "今日金价(元/克)",
				Data: template.Price.Round(2).String(),
			}},
		},
	}

	if _, err := wxwork.AgentWorkbench.SetWorkbenchTemplate(ctx, options); err != nil {
		log.Printf("SetWorkbenchTemplate error: %+v\n", err)
	}
}
