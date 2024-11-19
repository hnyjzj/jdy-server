package workbench

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/workbench"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type WorkbenchController struct {
	controller.BaseController
	logic workbench.WorkbenchLogic
}

// 获取列表
func (con WorkbenchController) List(ctx *gin.Context) {
	userinfo, err := con.logic.GetList()
	if err != nil {
		con.ErrorLogic(ctx, err)
		return
	}

	// 获取用户信息
	con.Success(ctx, "ok", userinfo)
}

// 添加入口
func (con WorkbenchController) Add(ctx *gin.Context) {

	var (
		req types.WorkbenchListReq
	)
	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := con.logic.AddRoute(&req)
	if err != nil {
		con.ErrorLogic(ctx, err)
		return
	}

	con.Success(ctx, "ok", res)
}
