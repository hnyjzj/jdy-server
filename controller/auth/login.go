package auth

import (
	"errors"
	"jdy/controller"
	authlogic "jdy/logic/auth"
	"jdy/logic_error"
	authtype "jdy/types/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	controller.BaseController
	logic authlogic.LoginLogic
}

func (con LoginController) Login(ctx *gin.Context) {
	var (
		req authtype.LoginReq
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, logic_error.ErrInvalidParam.Error())
		return
	}

	res, err := con.logic.Login(ctx, &req)
	if err != nil {
		// 验证码错误
		if errors.Is(err, logic_error.ErrInvalidCaptcha) {
			con.ErrorLogic(ctx, logic_error.ErrInvalidCaptcha)
			return
		}

		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

func (con LoginController) OAuth(ctx *gin.Context) {
	var (
		req authtype.LoginOAuthReq
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, logic_error.ErrInvalidParam.Error())
		return
	}

	res, err := con.logic.Oauth(ctx, &req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
