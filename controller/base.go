package controller

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"net/http"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 获取 token 中的用户信息
func (BaseController) GetStaff(ctx *gin.Context) *model.Staff {
	// 获取 token 中的用户信息
	staffInfo, ok := ctx.MustGet("staff").(*types.Staff)
	if staffInfo == nil || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    errors.ErrStaffNotFound.Code,
			"message": errors.ErrStaffNotFound.Message,
		})
		ctx.Abort()
		return nil
	}
	// 查询用户信息
	staff, db := gplus.SelectById[model.Staff](staffInfo.Id)
	if db.Error != nil || staff.Id == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    errors.ErrStaffNotFound.Code,
			"message": errors.ErrStaffNotFound.Message,
		})
		ctx.Abort()
		return nil
	}
	// 检查用户是否被禁用
	if staff.IsDisabled {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    errors.ErrStaffDisabled.Code,
			"message": errors.ErrStaffDisabled.Message,
		})
		ctx.Abort()
		return nil
	}

	return staff
}
