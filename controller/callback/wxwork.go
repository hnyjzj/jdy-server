package callback

import (
	"io"
	"jdy/config"
	"jdy/logic/callback"
	"net/http"

	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/gin-gonic/gin"
)

type WxWorkCongtroller struct {
	CallbackController
}

func (con WxWorkCongtroller) JdyVerify(c *gin.Context) {
	var (
		App = config.NewWechatService().JdyWork
	)

	con.verify(c, App)
}

func (con WxWorkCongtroller) ContactsVerify(c *gin.Context) {
	var (
		App = config.NewWechatService().ContactsWork
	)

	con.verify(c, App)
}

func (con WxWorkCongtroller) verify(c *gin.Context, App *work.Work) {
	rs, err := App.Server.VerifyURL(c.Request)
	if err != nil {
		panic(err)
	}
	text, _ := io.ReadAll(rs.Body)
	c.String(http.StatusOK, string(text))
}

const (
	TemplateCardEvent = "template_card_event" // 模板卡片事件
)

func (con WxWorkCongtroller) JdyNotify(c *gin.Context) {
	var (
		App = config.NewWechatService().JdyWork
	)

	con.notify(c, App)
}

func (con WxWorkCongtroller) ContactsNotify(c *gin.Context) {
	var (
		App = config.NewWechatService().ContactsWork
	)

	con.notify(c, App)
}

func (con WxWorkCongtroller) notify(c *gin.Context, App *work.Work) {

	var (
		logic = callback.WxWork{
			Ctx: c,
		}
	)

	rs, err := App.Server.Notify(c.Request, func(event contract.EventInterface) any {
		logic.Event = event
		var res any

		if event.GetEvent() == TemplateCardEvent {
			res = logic.TemplateCardEvent()
		}

		if res == nil {
			res = kernel.SUCCESS_EMPTY_RESPONSE
		}

		return res
	})
	if err != nil {
		panic(err)
	}

	err = helper.HttpResponseSend(rs, c.Writer)
	if err != nil {
		panic(err)
	}
}
