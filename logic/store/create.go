package store

import (
	"errors"
	"jdy/logic/platform/wxwork"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (l *StoreLogic) Create(ctx *gin.Context, req *types.StoreCreateReq) error {
	var (
		wxwork = wxwork.WxWorkLogic{}
	)

	var source_id int

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 创建企业微信部门
		if req.SyncWxwork {
			id, err := wxwork.StoreCreate(ctx, req)
			if err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
			source_id = id
		}

		if err := tx.Create(&model.Store{
			ParentId: req.ParentId,

			Name:     req.Name,
			Address:  req.Address,
			Contact:  req.Contact,
			Logo:     req.Logo,
			Order:    req.Order,
			Province: req.Province,
			City:     req.City,
			District: req.District,

			SourceId: source_id,
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		if req.SyncWxwork && source_id != 0 {
			if err := wxwork.StoreDelete(ctx, source_id); err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
		}
		return err
	}

	return nil
}
