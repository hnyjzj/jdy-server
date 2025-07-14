package setting

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/setting"
	"jdy/types"
	"jdy/utils"
	"log"

	"github.com/gin-gonic/gin"
)

type RemarkController struct {
	controller.BaseController
}

func (con RemarkController) Where(ctx *gin.Context) {
	where := utils.StructToWhere(types.RemarkWhere{})

	con.Success(ctx, "ok", where)
}

func (con RemarkController) Create(ctx *gin.Context) {

	var (
		req   types.RemarkCreateReq
		logic = &setting.RemarkLogic{}
	)

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	} else {
		// 设置上下文
		logic.Ctx = ctx
		logic.Staff = staff
	}

	if err := logic.Create(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con RemarkController) List(ctx *gin.Context) {

	var (
		req   types.RemarkListReq
		logic = &setting.RemarkLogic{}
	)

	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	} else {
		// 设置上下文
		logic.Ctx = ctx
		logic.Staff = staff
	}

	list, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

func (con RemarkController) Update(ctx *gin.Context) {
	var (
		req   types.RemarkUpdateReq
		logic = &setting.RemarkLogic{}
	)

	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	} else {
		// 设置上下文
		logic.Ctx = ctx
		logic.Staff = staff
	}

	if err := logic.Update(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}

func (con RemarkController) Delete(ctx *gin.Context) {
	var (
		req   types.RemarkDeleteReq
		logic = &setting.RemarkLogic{}
	)

	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("bind err: %v", err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	if staff, err := con.GetStaff(ctx); err != nil {
		con.Exception(ctx, err.Error())
		return
	} else {
		// 设置上下文
		logic.Ctx = ctx
		logic.Staff = staff
	}

	if err := logic.Delete(&req); err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", nil)
}
