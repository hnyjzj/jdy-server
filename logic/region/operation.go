package region

import (
	"errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (l *RegionLogic) Create(ctx *gin.Context, req *types.RegionCreateReq) error {
	region := &model.Region{
		Name: req.Name,
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(region).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (l *RegionLogic) Update(ctx *gin.Context, req *types.RegionUpdateReq) error {

	// 查询区域信息
	var region model.Region
	if err := model.DB.First(&region, "id = ?", req.Id).Error; err != nil {
		return errors.New("区域不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {

		data, err := utils.StructToStruct[model.Region](req)
		if err != nil {
			return errors.New("验证信息失败")
		}

		if err := tx.Model(&region).Where("id = ?", req.Id).Updates(data).Error; err != nil {
			return errors.New("更新失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// 删除区域
func (l *RegionLogic) Delete(ctx *gin.Context, req *types.RegionDeleteReq) error {
	// 查询区域信息
	region := &model.Region{}
	if err := model.DB.First(region, "id = ?", req.Id).Error; err != nil {
		return errors.New("区域不存在或已被删除")
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(region).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
