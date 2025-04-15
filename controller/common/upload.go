package common

import (
	"jdy/controller"
	"jdy/errors"
	"jdy/logic/common"
	"jdy/types"
	"log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	controller.BaseController
}

// 上传头像
func (con UploadController) Avatar(ctx *gin.Context) {
	// 接收参数
	type Req struct {
		File *multipart.FileHeader `form:"avatar" binding:"required"`
	}
	type Res struct {
		Url string `json:"url"`
	}
	var (
		r Req
		s Res
	)

	staff, err := con.GetStaff(ctx)
	if err != nil {
		con.ExceptionWithAuth(ctx, err.Error())
		return
	}

	// 验证参数
	if err := ctx.ShouldBind(&r); err != nil {
		log.Println(err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 上传文件
	upload := &common.Upload{
		Ctx:    ctx,
		File:   r.File,
		Model:  types.UploadModelAvatar,
		Type:   types.UploadTypeImage,
		Prefix: staff.Id,
	}

	url, err := upload.Save()
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	s.Url = url.Uris[0]
	con.Success(ctx, "ok", s)
}

// 上传工作台图标
func (con UploadController) Workbench(ctx *gin.Context) {
	// 接收参数
	type Req struct {
		File *multipart.FileHeader `form:"image" binding:"required"`
	}
	type Res struct {
		Url string `json:"url"`
	}
	var (
		r Req
		s Res
	)

	// 验证参数
	if err := ctx.ShouldBind(&r); err != nil {
		log.Println(err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 上传文件
	upload := &common.Upload{
		Ctx:   ctx,
		File:  r.File,
		Model: types.UploadModelWorkbench,
		Type:  types.UploadTypeImage,
	}
	url, err := upload.Save()
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	s.Url = url.Uris[0]
	con.Success(ctx, "ok", s)
}

// 上传门店图标
func (con UploadController) Store(ctx *gin.Context) {
	// 接收参数
	type Req struct {
		File    *multipart.FileHeader `form:"image" binding:"required"`
		StoreId string                `form:"store_id"`
	}
	type Res struct {
		Url string `json:"url"`
	}
	var (
		r Req
		s Res
	)

	// 验证参数
	if err := ctx.ShouldBind(&r); err != nil {
		log.Println(err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 上传文件
	upload := &common.Upload{
		Ctx:    ctx,
		File:   r.File,
		Model:  types.UploadModelStore,
		Type:   types.UploadTypeImage,
		Prefix: r.StoreId,
	}
	url, err := upload.Save()
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	s.Url = url.Uris[0]
	con.Success(ctx, "ok", s)
}

// 上传商品图片
func (con UploadController) Product(ctx *gin.Context) {
	// 接收参数
	type Req struct {
		File      *multipart.FileHeader `form:"image" binding:"required"`
		ProductId string                `form:"product_id" binding:"required"`
	}
	type Res struct {
		Url string `json:"url"`
	}
	var (
		r Req
		s Res
	)
	// 验证参数
	if err := ctx.ShouldBind(&r); err != nil {
		log.Println(err)
		con.Exception(ctx, errors.ErrInvalidParam.Error())
		return
	}

	// 上传文件
	upload := &common.Upload{
		Ctx:    ctx,
		File:   r.File,
		Model:  types.UploadModelProduct,
		Type:   types.UploadTypeImage,
		Prefix: r.ProductId,
	}
	url, err := upload.Save()
	if err != nil {
		con.Exception(ctx, err.Error())
		return
	}

	s.Url = url.Uris[0]
	con.Success(ctx, "ok", s)
}
