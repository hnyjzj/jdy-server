package auth

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/auth"
	"jdy/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	controller.BaseController
}

func (con LoginController) Login(ctx *gin.Context) {
	var (
		req   types.LoginReq
		logic = auth.LoginLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.Login(ctx, &req)
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
		req   types.LoginOAuthReq
		logic = auth.LoginLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.Oauth(ctx, &req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}

func (con LoginController) Logout(ctx *gin.Context) {
	var (
		logic = auth.LoginLogic{}
	)

	staff := con.GetStaff(ctx)

	if err := logic.Logout(ctx, staff.Phone); err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
