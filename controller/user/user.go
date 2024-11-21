package user

import (
	"fmt"
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/user"
	"jdy/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	controller.BaseController
}

// 创建用户
func (con UserController) Create(ctx *gin.Context) {
	var (
		req types.UserReq

		logic = user.UserLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 创建用户
	user, err := logic.CreateUser(ctx, &req)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回结果
	con.Success(ctx, "ok", user)
}

// 获取用户信息
func (con UserController) Info(ctx *gin.Context) {
	var (
		logic = user.UserLogic{}
	)
	user := con.GetUser(ctx)

	userinfo, err := logic.GetUserInfo(user.Id)
	if err != nil {
		con.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取用户信息
	con.Success(ctx, "ok", userinfo)
}
