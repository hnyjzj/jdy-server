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
			Ctx:   ctx,
			Staff: con.GetStaff(ctx),
		}
	)

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
