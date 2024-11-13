package user

import (
	"jdy/controller"
	userlogic "jdy/logic/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	controller.BaseController
}

func (con UserController) Info(ctx *gin.Context) {
	var (
		logic userlogic.UserLogic
	)
	user := con.GetUser(ctx)

	userinfo, err := logic.GetUserInfo(user.Id)
	if err != nil {
		con.ErrorJson(ctx, http.StatusNotFound, "用户不存在")
		return
	}

	// 获取用户信息
	con.SuccessJson(ctx, "ok", userinfo)
}
