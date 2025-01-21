package callback

import (
	"fmt"
	"io"
	"jdy/config"
	"net/http"

	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"github.com/gin-gonic/gin"
)

type WxWorkCongtroller struct {
	CallbackController
}

func (con WxWorkCongtroller) JdyVerify(c *gin.Context) {
	var (
		Jwt = config.NewWechatService().JdyWork
	)

	rs, err := Jwt.Server.Serve(c.Request)
	fmt.Printf("err.Error(): %v\n", err.Error())
	// if err != nil {
	// 	panic(err.Error())
	// }

	text, _ := io.ReadAll(rs.Body)
	c.String(http.StatusOK, string(text))
}

func (con WxWorkCongtroller) JdyNotify(c *gin.Context) {

	var (
		Jwt = config.NewWechatService().JdyWork
	)
	rs, err := Jwt.Server.Notify(c.Request, func(event contract.EventInterface) any {
		// 这里需要获取到事件类型，然后把对应的结构体传递进去进一步解析
		// 所有包含的结构体请参考： https://github.com/ArtisanCloud/PowerWeChat/tree/master/src/work/server/handlers/models
		if event.GetEvent() == models.CALLBACK_EVENT_CHANGE_CONTACT && event.GetChangeType() == models.CALLBACK_EVENT_CHANGE_TYPE_CREATE_PARTY {
			//msg := models.EventPartyCreate{}
			msg := models.EventKFMsgOrEvent{}
			err := event.ReadMessage(&msg)
			if err != nil {
				println(err.Error())
				return "error"
			}
		}

		// 假设员工给应用发送消息，这里可以直接回复消息文本，
		// return  "I'm recv..."

		// 这里回复success告诉微信我收到了，后续需要回复用户信息可以主动调发消息接口
		return kernel.SUCCESS_EMPTY_RESPONSE
	})
	if err != nil {
		panic(err)
	}

	err = helper.HttpResponseSend(rs, c.Writer)
	if err != nil {
		panic(err)
	}
}
