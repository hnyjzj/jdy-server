package auth

import (
	"jdy/controller"
	"jdy/errors"
	authlogic "jdy/logic/auth"
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
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := con.logic.Login(ctx, &req)
	if err != nil {
		// 验证码错误
		if errors.Is(err, errors.ErrInvalidCaptcha) {
			con.ErrorLogic(ctx, errors.ErrInvalidCaptcha)
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
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := con.logic.Oauth(ctx, &req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
