package auth

import (
	"fmt"
	"jdy/controller"
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
		fmt.Printf("Login err: %v \n", err.Error())
		con.ExceptionJson(ctx, "参数错误")
		return
	}

	res, err := con.logic.Login(&req)
	if err != nil {
		con.ErrorJson(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.SuccessJson(ctx, "ok", res)
}

func (con LoginController) OAuth(ctx *gin.Context) {
	var (
		req authtype.LoginOAuthReq
	)
	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.ExceptionJson(ctx, "参数错误")
		return
	}

	res, err := con.logic.Oauth(&req)
	if err != nil {
		con.ErrorJson(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	con.SuccessJson(ctx, "ok", res)
}
