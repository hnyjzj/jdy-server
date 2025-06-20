package controller

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"log"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 获取 token 中的用户信息
func (con BaseController) GetStaff(ctx *gin.Context) (*model.Staff, *errors.Errors) {
	// 获取 token 中的用户信息
	staffInfo, ok := ctx.MustGet("staff").(*types.Staff)
	// 检查用户是否正确
	if !ok || staffInfo == nil {
		return nil, errors.ErrStaffNotFound
	}

	staff, err := model.Staff{}.Get(staffInfo.Id)
	if err != nil {
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
		log.Printf("员工[%v] 无权限访问: %v", staff.Id, ctx.FullPath())
		return errors.ErrStaffUnauthorized
	}

	return nil
}
