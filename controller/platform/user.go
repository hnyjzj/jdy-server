package platform

import (
	"jdy/errors"
	"jdy/logic/platform"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

// 三方 JSSDK
func (con PlatformController) GetUser(ctx *gin.Context) {

	// 绑定参数
	var (
		req   types.PlatformGetUserReq
		logic = platform.PlatformLogic{
			Ctx: ctx,
		}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.GetUser(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
