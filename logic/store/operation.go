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

	store := &model.Store{
		ParentId: req.ParentId,

		Name:     req.Name,
		Address:  req.Address,
		Contact:  req.Contact,
		Logo:     req.Logo,
		Sort:     req.Sort,
		Province: req.Province,
		City:     req.City,
		District: req.District,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 创建企业微信部门
		if req.SyncWxwork {
			id, err := wxwork.StoreCreate(ctx, req)
			if err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
			store.WxworkId = id
		}

		if err := tx.Create(store).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		if req.SyncWxwork && store.WxworkId != 0 {
			if err := wxwork.StoreDelete(ctx, store.WxworkId); err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
		}
		return err
	}

	return nil
}

func (l *StoreLogic) Update(ctx *gin.Context, req *types.StoreUpdateReq) error {
	var (
		wxwork = wxwork.WxWorkLogic{}
	)

	// 查询门店信息
	var store model.Store
	if err := model.DB.First(store, req.Id).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 更新企业微信部门
		if req.SyncWxwork && store.WxworkId != 0 {
			if err := wxwork.StoreDelete(ctx, store.WxworkId); err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
		}

		if err := tx.Save(req).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		if req.SyncWxwork && store.WxworkId != 0 {
			if err := wxwork.StoreDelete(ctx, store.WxworkId); err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
		}
		return err
	}

	return nil
}

// 删除门店
func (l *StoreLogic) Delete(ctx *gin.Context, req *types.StoreDeleteReq) error {
	var (
		wxwork = wxwork.WxWorkLogic{}
	)

	// 查询门店信息
	store := &model.Store{}
	if err := model.DB.First(store, req.Id).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 删除企业微信部门
		if req.SyncWxwork && store.WxworkId != 0 {
			if err := wxwork.StoreDelete(ctx, store.WxworkId); err != nil {
				return errors.New("同步企业微信失败: " + err.Error())
			}
		}

		if err := tx.Delete(store).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
