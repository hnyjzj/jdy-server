package controller

import (
	"jdy/enums"
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

	staff, err := model.Staff{}.Get(&staffInfo.Id, nil)
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

	if err := con.Verify_permission(ctx, staff); err != nil {
		return nil, &errors.Errors{
			Message: err.Error(),
			Code:    errors.ErrStaffUnauthorized.Code,
		}
	}

	return staff, nil
}

func (con BaseController) GetLoginType(ctx *gin.Context) (enums.LoginType, *errors.Errors) {
	// 获取 token 中的用户信息
	ltype, ok := ctx.MustGet("type").(enums.LoginType)
	// 检查用户是否正确
	if !ok {
		return "", &errors.Errors{
			Message: "登录错误",
			Code:    errors.ErrStaffUnauthorized.Code,
		}
	}

	return ltype, nil
}

func (con BaseController) Verify_permission(ctx *gin.Context, staff *model.Staff) error {
	if staff.Identity == enums.IdentitySuperAdmin {
		return nil
	}

	// 检查权限
	if err := staff.HasPermissionApi(ctx.FullPath()); err != nil {
		log.Printf("%s[%v] %v", staff.Nickname, staff.Id, err.Error())
		return err
	}

	return nil
}
