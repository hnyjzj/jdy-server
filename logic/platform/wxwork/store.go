package wxwork

import (
	"errors"
	"jdy/config"
	"jdy/types"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/department/request"
	"github.com/gin-gonic/gin"
)

func (w *WxWorkLogic) StoreCreate(ctx *gin.Context, req *types.StoreCreateReq) (id int, err error) {
	wxwork := config.NewWechatService().ContactsWork

	params := &request.RequestDepartmentInsert{
		Name:     req.Name,
		Order:    1,
		ParentID: 1,
	}
	if req.Sort > 0 {
		params.Order = req.Sort
	}

	res, err := wxwork.Department.Create(ctx, params)
	if err != nil || res.ErrCode != 0 {
		return 0, errors.New(res.ErrMsg)
	}

	id = res.ID

	return id, nil
}

func (w *WxWorkLogic) StoreUpdate(ctx *gin.Context, id int, req *types.StoreUpdateReq) error {
	wxwork := config.NewWechatService().ContactsWork

	params := &request.RequestDepartmentUpdate{
		ID:       id,
		Name:     req.Name,
		Order:    1,
		ParentID: 1,
	}
	if req.Sort > 0 {
		params.Order = req.Sort
	}
	res, err := wxwork.Department.Update(ctx, params)
	if err != nil || res.ErrCode != 0 {
		return errors.New(res.ErrMsg)
	}

	return nil
}

func (w *WxWorkLogic) StoreDelete(ctx *gin.Context, id int) error {
	wxwork := config.NewWechatService().ContactsWork

	res, err := wxwork.Department.Delete(ctx, id)
	if err != nil || res.ErrCode != 0 {
		return errors.New(res.ErrMsg)
	}

	return nil
}
