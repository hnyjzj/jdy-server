package common

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/message"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	controller.BaseController
}

func (con LogController) OnCaptureScreen(ctx *gin.Context) {
	var (
		req types.OnCaptureScreenReq
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 发送消息
	m := message.NewMessage(ctx)
	if err := m.SendCaptureScreenMessage(&message.CaptureScreenMessage{
		Username:  req.Username,
		Storename: req.Storename,
		Url:       req.Url,
	}); err != nil {
		con.Exception(ctx, "记录失败")
		return
	}

	con.Success(ctx, "记录成功", nil)
}
