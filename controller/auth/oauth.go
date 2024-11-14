package auth

import (
	"jdy/controller"
	authlogic "jdy/logic/auth"
	"jdy/logic_error"
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
		authlogic authlogic.OAuthLogic
	)

	// 获取请求头
	req.Agent = ctx.GetHeader("User-Agent")

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, logic_error.ErrInvalidParam.Error())
		return
	}

	res, err := authlogic.GetUri(&req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
