package platform

import (
	"jdy/errors"
	"jdy/logic/platform"
	"jdy/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 三方授权链接
func (con PlatformController) OauthUri(ctx *gin.Context) {
	// 绑定参数
	var (
		req   types.PlatformOAuthReq
		logic = platform.PlatformLogic{}
	)

	// 获取请求头
	req.Agent = ctx.GetHeader("User-Agent")

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.OauthUri(&req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
