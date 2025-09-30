package store

import (
	"errors"
	"jdy/enums"
	"jdy/logic/platform"
	"jdy/model"
	"jdy/types"
	"jdy/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (l *StoreLogic) Create(ctx *gin.Context, req *types.StoreCreateReq) error {
	// 验证门店名称是否以"门店"结尾
	if !strings.HasSuffix(req.Name, enums.DepartmentStore.String()) {
		return errors.New("门店名称必须以" + enums.DepartmentStore.String() + "结尾")
	}

	// 转换结构体
	store, err := utils.StructToStruct[model.Store](req)
	if err != nil {
		return errors.New("验证信息失败")
	}

	// 查询区域
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.RegionId).Error; err != nil {
		return errors.New("区域不存在或已被删除")
	}

	// 创建门店
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 同步企业微信
		pf := &platform.DepartmentLogic{}
		pf.Ctx = ctx

		regionId, err := strconv.Atoi(region.IdWx)
		if err != nil {
			return errors.New("区域ID转换失败")
		}
		id, err := pf.Create(&types.PlatformDepartmentCreateReq{
			Name:     req.Name,
			NameEn:   req.Alias,
			ParentId: regionId,
			Order:    0,
		})
		if err != nil {
			return errors.New("同步企业微信失败")
		}

		store.IdWx = strconv.Itoa(id)

		if err := tx.Create(&store).Error; err != nil {
			return errors.New("创建失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (l *StoreLogic) Update(ctx *gin.Context, req *types.StoreUpdateReq) error {

	// 查询门店信息
	var store model.Store
	if err := model.DB.First(&store, "id = ?", req.Id).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		data, err := utils.StructToStruct[model.Store](req)
		if err != nil {
			return errors.New("验证信息失败")
		}

		if err := tx.Model(&model.Store{}).Where("id = ?", store.Id).Updates(data).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 删除门店
func (l *StoreLogic) Delete(ctx *gin.Context, req *types.StoreDeleteReq) error {
	// 查询门店信息
	store := &model.Store{}
	if err := model.DB.First(store, "id = ?", req.Id).Error; err != nil {
		return errors.New("门店不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(store).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
