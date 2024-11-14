package user

import (
	"jdy/controller"
	"jdy/errors"
	userlogic "jdy/logic/user"

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
		con.ErrorLogic(ctx, errors.ErrUserNotFound)
		return
	}

	// 获取用户信息
	con.Success(ctx, "ok", userinfo)
}
