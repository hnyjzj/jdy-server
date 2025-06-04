package callback

import (
	"io"
	"jdy/config"
	"jdy/logic/callback"
	"log"
	"net/http"

	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/gin-gonic/gin"
)

const (
	EventTemplateCard          = "template_card_event"     // 模板卡片事件
	EventChangeExternalContact = "change_external_contact" // 客户变更事件
	EventChangeContact         = "change_contact"          // 通讯录变更事件
)

func (con WxWorkCongtroller) notify(c *gin.Context, App *work.Work) {

	handler := callback.WxWork{
		Ctx: c,
	}

	rs, err := App.Server.Notify(c.Request, func(event contract.EventInterface) any {
		handler.Event = event
		var res any

		switch event.GetEvent() {
		case EventTemplateCard:
			res = handler.TemplateCardEvent()
		case EventChangeExternalContact:
			res = handler.ChangeExternalContactEvent()
		case EventChangeContact:
			res = handler.ChangeContactEvent()
		default:
			log.Printf("wxwork notify event: %s", event.GetEvent())
		}

		if res == nil {
			res = kernel.SUCCESS_EMPTY_RESPONSE
		}

		return res
	})

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("wxwork notify error: %+v", err)
		return
	}

	err = helper.HttpResponseSend(rs, c.Writer)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("wxwork notify error: %+v", err)
		return
	}
}

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
