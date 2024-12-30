package member

import "github.com/gin-gonic/gin"

func (con MemberController) Create(ctx *gin.Context) {
	con.Success(ctx, "ok", nil)
}

func (con MemberController) Update(ctx *gin.Context) {
	con.Success(ctx, "ok", nil)
}
