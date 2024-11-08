package auth

import (
	"jdy/controller"
	authlogic "jdy/logic/auth"
	authtype "jdy/types/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthController struct {
	controller.BaseController
}

// 企业微信授权登录
func (con OAuthController) GetUri(ctx *gin.Context) {
	// 绑定参数
	var (
		req       authtype.OAuthWeChatWorkReq
		authlogic = authlogic.OAuthLogic{}
	)

	// 获取请求头
	req.Agent = ctx.GetHeader("User-Agent")

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.ExceptionJson(ctx, "参数错误")
		return
	}

	res, err := authlogic.GetUri(&req)
	if err != nil {
		con.ErrorJson(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.SuccessJson(ctx, "ok", res)
}
