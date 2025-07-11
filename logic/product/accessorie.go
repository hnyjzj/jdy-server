package product

import (
	"jdy/enums"
	"jdy/errors"
	"jdy/model"
	"jdy/types"
	"jdy/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductAccessorieLogic struct {
	Ctx   *gin.Context
	Staff *model.Staff
}

// 配件列表
func (p *ProductAccessorieLogic) List(req *types.ProductAccessorieListReq) (*types.PageRes[model.ProductAccessorie], error) {
	var (
		product model.ProductAccessorie

		res types.PageRes[model.ProductAccessorie]
	)

	db := model.DB.Model(&product)
	db = product.WhereCondition(db, &req.Where)

	// 获取总数
	if err := db.Count(&res.Total).Error; err != nil {
		return nil, errors.New("获取配件数量失败")
	}

	// 获取列表
	db = model.PageCondition(db, req.Page, req.Limit)
	db = db.Order("created_at desc")
	db = db.Preload("Category")
	// db = db.Select("*,SUM(stock) as stock")
	// group := []string{
	// 	"id", "created_at", "updated_at", "deleted_at",
	// 	"store_id",
	// 	"code",
	// 	"stock",
	// 	"access_fee",
	// 	"status",
	// 	"enter_id",
	// }
	// db = db.Group(strings.Join(group, ","))
	if err := db.Find(&res.List).Error; err != nil {
		return nil, errors.New("获取配件列表失败")
	}

	return &res, nil
}

// 配件详情
func (p *ProductAccessorieLogic) Info(req *types.ProductAccessorieInfoReq) (*model.ProductAccessorie, error) {
	var (
		product model.ProductAccessorie
	)

	db := model.DB.Model(&model.ProductAccessorie{})
	db = db.Where("id = ?", req.Id)
	db = db.Preload("Category")
	db = db.Preload("Store")

	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("获取配件信息失败")
	}

	return &product, nil
}

// 更新配件信息
func (p *ProductAccessorieLogic) Update(req *types.ProductAccessorieUpdateReq) error {

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		data, err := utils.StructToStruct[model.ProductAccessorie](req)
		if err != nil {
			return errors.New("验证参数失败")
		}

		var product model.ProductAccessorie
		if err := tx.Model(&model.ProductAccessorie{}).
			Preload("Store").
			Where("id = ?", req.Id).First(&product).Error; err != nil {
			return errors.New("获取配件信息失败")
		}

		history := model.ProductHistory{
			Type:       enums.ProductTypeAccessorie,
			Action:     enums.ProductActionUpdate,
			OldValue:   product,
			ProductId:  product.Id,
			SourceId:   product.Id,
			StoreId:    product.StoreId,
			OperatorId: p.Staff.Id,
			IP:         p.Ctx.ClientIP(),
		}

		if err := tx.Model(&product).Clauses(clause.Returning{}).Where("id = ?", req.Id).Updates(&data).Error; err != nil {
			return errors.New("更新配件信息失败")
		}

		// 添加记录
		history.NewValue = product
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
