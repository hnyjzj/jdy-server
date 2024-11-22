package platform

import (
	"fmt"
	"jdy/errors"
	"jdy/logic/platform"
	"jdy/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 三方 JSSDK
func (con PlatformController) JSSDK(ctx *gin.Context) {

	// 绑定参数
	var (
		req   types.PlatformJSSdkReq
		logic = platform.PlatformLogic{}
	)

	// 获取请求头
	req.Agent = ctx.GetHeader("User-Agent")

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.GetJSSDK(ctx, &req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", res)

}
