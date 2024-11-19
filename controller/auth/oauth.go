package auth

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/auth"
	"jdy/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthController struct {
	controller.BaseController
}

// 三方授权链接
func (con OAuthController) GetOauthUri(ctx *gin.Context) {
	// 绑定参数
	var (
		req       types.OAuthWeChatWorkReq
		authlogic auth.OAuthLogic
	)

	// 获取请求头
	req.Agent = ctx.GetHeader("User-Agent")

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := authlogic.OauthUri(&req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
