package controller

import (
	"jdy/logic_error"
	usermodel "jdy/model/user"
	"net/http"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 获取 token 中的用户信息
func (BaseController) GetUser(ctx *gin.Context) *usermodel.User {
	userInfo := ctx.MustGet("user").(usermodel.User)
	user, db := gplus.SelectById[usermodel.User](userInfo.Id)

	if db.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    logic_error.ErrUserNotFound.Code,
			"message": logic_error.ErrUserNotFound.Message,
		})
		ctx.Abort()
		return nil
	}

	if user.IsDisabled {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    logic_error.ErrUserDisabled.Code,
			"message": logic_error.ErrUserDisabled.Message,
		})
		ctx.Abort()
		return nil
	}

	return user
}
