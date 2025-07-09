package wxwork

import (
	"errors"
	"jdy/config"
	"jdy/types"
	"log"

	"github.com/gin-gonic/gin"
)

type PlatformJSSdkRes struct {
	AgentTicket string `json:"agent_ticket"` // jsapi_ticket
	JsapiTicket string `json:"jsapi_ticket"` // jsapi_ticket
	CorpID      string `json:"corp_id"`      // 企业ID
	AgentID     int    `json:"agent_id"`     // 应用ID
}

func (w *WxWorkLogic) Jssdk(ctx *gin.Context, req *types.PlatformJSSdkReq) (*PlatformJSSdkRes, error) {
	var (
		wxwork = config.NewWechatService().JdyWork
	)

	config := wxwork.GetConfig()

	agent, err := wxwork.JSSDK.GetTicket(ctx)
	if err != nil || (agent != nil && agent.ErrCode != 0) {
		log.Printf("获取 ticket 失败, err: %v, agent: %+v", err, agent)
		return nil, errors.New("获取 ticket 失败")
	}

	res := &PlatformJSSdkRes{
		AgentTicket: agent.Ticket,
		CorpID:      config.GetString("corp_id", ""),
		AgentID:     config.GetInt("agent_id", 0),
	}

	return res, nil
}
