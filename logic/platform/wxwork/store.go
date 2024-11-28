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

	res, err := wxwork.Department.Create(ctx, &request.RequestDepartmentInsert{
		Name:     req.Name,
		ParentID: req.SourceId,
		Order:    req.Order,
	})
	if err != nil || res.ErrCode != 0 {
		return 0, errors.New(res.ErrMsg)
	}

	id = res.ID

	return id, nil
}

func (w *WxWorkLogic) StoreDelete(ctx *gin.Context, id int) error {
	wxwork := config.NewWechatService().ContactsWork

	res, err := wxwork.Department.Delete(ctx, id)
	if err != nil || res.ErrCode != 0 {
		return errors.New(res.ErrMsg)
	}

	return nil
}
