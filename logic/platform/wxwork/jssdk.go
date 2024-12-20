package wxwork

import (
	"errors"
	"jdy/config"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

func (w *WxWorkLogic) Jssdk(ctx *gin.Context, req *types.PlatformJSSdkReq) (*string, error) {
	var (
		wxwork = config.NewWechatService().JdyWork
		client = wxwork.JSSDK.Client
	)

	wxwork.GetAccessToken()
	switch req.Type {
	case "jsapi":
		client.TicketEndpoint = "/cgi-bin/get_jsapi_ticket"
		jsapi, err := wxwork.JSSDK.Client.GetTicket(ctx, false, "jsapi")
		if err != nil {
			return nil, errors.New("jsapi ticket error")
		}

		ticket, ok := jsapi.Get("ticket").(string)
		if !ok {
			return nil, errors.New("jsapi ticket error")
		}

		return &ticket, err

	case "agent":
		client.TicketEndpoint = "/cgi-bin/ticket/get"
		agent, err := wxwork.JSSDK.Client.GetTicket(ctx, false, "agent_config")
		if err != nil {
			return nil, errors.New("agent ticket error")
		}

		ticket, ok := agent.Get("ticket").(string)
		if !ok {
			return nil, errors.New("agent ticket error")
		}

		return &ticket, err

	default:
		return nil, errors.New("state error")
	}
}
