package platform

import (
	"errors"
	"jdy/config"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 获取授权链接
func (l *PlatformLogic) GetJSSDK(ctx *gin.Context, req *types.PlatformJSSdkReq) (*types.PlatformJSSdkRes, error) {
	switch req.Platform {
	case "wxwork_jsapi":
		var (
			wxwork = config.NewWechatService().JdyWork
			client = wxwork.JSSDK.Client
		)
		wxwork.GetAccessToken()
		client.TicketEndpoint = "/cgi-bin/get_jsapi_ticket"
		jsapi, err := wxwork.JSSDK.Client.GetTicket(ctx, false, "jsapi")
		if err != nil {
			return nil, errors.New("jsapi ticket error")
		}

		res := types.PlatformJSSdkRes{
			Ticket: jsapi.Get("ticket").(string),
		}

		return &res, err

	case "wxwork_agent":
		var (
			wxwork = config.NewWechatService().JdyWork
			client = wxwork.JSSDK.Client
		)
		wxwork.GetAccessToken()

		client.TicketEndpoint = "/cgi-bin/ticket/get"
		agent, err := wxwork.JSSDK.Client.GetTicket(ctx, false, "agent_config")
		if err != nil {
			return nil, errors.New("agent ticket error")
		}

		res := types.PlatformJSSdkRes{
			Ticket: agent.Get("ticket").(string),
		}

		return &res, err

	default:
		return nil, errors.New("state error")
	}

}
