package product

import (
	"jdy/errors"
	"jdy/model"
	"jdy/types"

	"github.com/gin-gonic/gin"
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
	db = model.PageCondition(db, &req.PageReq)
	db = db.Order("created_at desc")
	db = db.Order("stock desc")
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
	db = product.Preloads(db)

	if err := db.First(&product).Error; err != nil {
		return nil, errors.New("获取配件信息失败")
	}

	return &product, nil
}
