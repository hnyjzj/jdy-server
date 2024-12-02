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
}

// 获取列表
func (con WorkbenchController) List(ctx *gin.Context) {
	var (
		logic = workbench.WorkbenchLogic{}
	)
	workbenchs, err := logic.GetList()
	if err != nil {
		con.ErrorLogic(ctx, err)
		return
	}

	// 获取用户信息
	con.Success(ctx, "ok", workbenchs)
}

// 添加入口
func (con WorkbenchController) Add(ctx *gin.Context) {
	var (
		req types.WorkbenchAddReq

		logic = workbench.WorkbenchLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.AddRoute(&req)
	if err != nil {
		con.ErrorLogic(ctx, err)
		return
	}

	con.Success(ctx, "ok", res)
}

// 删除入口
func (con WorkbenchController) Del(ctx *gin.Context) {
	var (
		req types.WorkbenchDelReq

		logic = workbench.WorkbenchLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.DelRoute(req.Id); err != nil {
		con.ErrorLogic(ctx, err)
		return
	}

	con.Success(ctx, "ok", nil)
}

// 更新入口
func (con WorkbenchController) Update(ctx *gin.Context) {
	var (
		req types.WorkbenchUpdateReq

		logic = workbench.WorkbenchLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if err := logic.UpdateRoute(&req); err != nil {
		con.ErrorLogic(ctx, err)
		return
	}

	con.Success(ctx, "ok", nil)
}
