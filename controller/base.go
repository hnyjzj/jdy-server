package controller

import (
	"jdy/errors"
	"jdy/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 获取 token 中的用户信息
func (con BaseController) GetStaff(ctx *gin.Context) *types.Staff {
	// 获取 token 中的用户信息
	staffInfo, ok := ctx.MustGet("staff").(*types.Staff)
	// 检查用户是否正确
	if staffInfo == nil || !ok {
		con.Error(ctx, http.StatusUnauthorized, errors.ErrStaffNotFound.Error())
		return nil
	}

	// 判断 IP
	if staffInfo.IP != ctx.ClientIP() {
		con.Error(ctx, http.StatusUnauthorized, "IP 地址不匹配")
		return nil
	}

	if staffInfo.UserAgent != ctx.GetHeader("User-Agent") {
		con.Error(ctx, http.StatusUnauthorized, "User-Agent 不匹配")
		return nil
	}

	// 检查用户是否被禁用
	if staffInfo.IsDisabled {
		con.Error(ctx, http.StatusUnauthorized, errors.ErrStaffDisabled.Error())
		return nil
	}

	return staffInfo
}
