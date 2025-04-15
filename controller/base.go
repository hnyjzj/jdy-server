package controller

import (
	"jdy/errors"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 获取 token 中的用户信息
func (con BaseController) GetStaff(ctx *gin.Context) (*types.Staff, error) {
	// 获取 token 中的用户信息
	staffInfo, ok := ctx.MustGet("staff").(*types.Staff)

	// 检查用户是否正确
	if staffInfo == nil || !ok {
		return nil, errors.ErrStaffNotFound
	}

	// 判断 IP
	// if staffInfo.IP != ctx.ClientIP() {
	// 	return nil, errors.New("IP 地址不匹配")
	// }

	// 检查用户是否被禁用
	if staffInfo.IsDisabled {
		return nil, errors.ErrStaffDisabled
	}

	return staffInfo, nil
}
