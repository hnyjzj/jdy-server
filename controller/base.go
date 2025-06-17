package controller

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 获取 token 中的用户信息
func (con BaseController) GetStaff(ctx *gin.Context) (*model.Staff, error) {
	// 获取 token 中的用户信息
	staffInfo, ok := ctx.MustGet("staff").(*types.Staff)

	staff, err := model.Staff{}.Get(staffInfo.Id)
	if err != nil {
		return nil, err
	}

	// 检查用户是否正确
	if staff == nil || !ok {
		return nil, errors.ErrStaffNotFound
	}

	// 判断 IP
	// if staff.IP != ctx.ClientIP() {
	// 	return nil, errors.New("IP 地址不匹配")
	// }

	// 检查用户是否被禁用
	if staff.IsDisabled {
		return nil, errors.ErrStaffDisabled
	}

	if err := con.verify_permission(ctx, staff); err != nil {
		return nil, errors.ErrStaffUnauthorized
	}

	return staff, nil
}

func (con BaseController) verify_permission(ctx *gin.Context, staff *model.Staff) error {
	if staff.IsRoot() {
		return nil
	}

	if !staff.HasPermissionApi(ctx.FullPath()) {
		return errors.ErrStaffUnauthorized
	}

	return nil
}
