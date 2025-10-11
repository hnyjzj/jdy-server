package store

import (
	"jdy/errors"
	"jdy/logic/store"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreStaffController struct {
	StoreController
}

// 门店员工列表
func (con StoreStaffController) List(ctx *gin.Context) {
	var (
		req   types.StoreStaffListReq
		logic = store.StoreStaffLogic{
			Ctx: ctx,
		}
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 查询门店列表
	list, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

func (con StoreStaffController) IsIn(ctx *gin.Context) {
	var (
		req   types.StoreStaffIsInReq
		logic = store.StoreStaffLogic{
			Ctx: ctx,
		}
	)

	if staff, err := con.GetStaff(ctx); err != nil {
		con.ExceptionWithAuth(ctx, err)
		return
	} else {
		logic.Staff = staff
	}

	// 校验参数
	if err := ctx.ShouldBind(&req); err != nil {
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	res, err := logic.IsIn(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
