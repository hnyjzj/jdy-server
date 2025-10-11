package store

import (
	"jdy/errors"
	"jdy/logic/store"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type StoreSuperiorController struct {
	StoreController
}

// 门店负责人列表
func (con StoreSuperiorController) List(ctx *gin.Context) {
	var (
		req   types.StoreSuperiorListReq
		logic = store.StoreSuperiorLogic{
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

	// 查询负责人列表
	list, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}

func (con StoreSuperiorController) IsIn(ctx *gin.Context) {
	var (
		req   types.StoreSuperiorIsInReq
		logic = store.StoreSuperiorLogic{
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

	// 查询
	res, err := logic.IsIn(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", res)
}
