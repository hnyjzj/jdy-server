package region

import (
	"jdy/errors"
	"jdy/logic/region"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type RegionStaffController struct {
	RegionController
}

// 区域员工列表
func (con RegionStaffController) List(ctx *gin.Context) {
	var (
		req   types.RegionStaffListReq
		logic = region.RegionStaffLogic{
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

	// 查询区域列表
	list, err := logic.List(&req)
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	con.Success(ctx, "ok", list)
}
