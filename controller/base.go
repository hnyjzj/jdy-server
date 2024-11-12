package controller

import (
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
			"code":    http.StatusUnauthorized,
			"message": "用户不存在",
		})
		ctx.Abort()
		return nil
	}

	if user.IsDisabled {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "用户已被禁用",
		})
		ctx.Abort()
		return nil
	}

	return user
}
