package controller

import (
	"jdy/errors"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 获取 token 中的用户信息
func (con BaseController) GetStaff(ctx *gin.Context) *types.Staff {
	// 获取 token 中的用户信息
	staffInfo, ok := ctx.MustGet("staff").(*types.Staff)
	// 检查用户是否正确
	if staffInfo == nil || !ok {
		con.Exception(ctx, errors.ErrStaffNotFound.Error())
		return nil
	}
	// 检查用户是否被禁用
	if staffInfo.IsDisabled {
		con.Exception(ctx, errors.ErrStaffDisabled.Error())
		return nil
	}

	return staffInfo
}
