package statistic

import (
	"errors"
	"jdy/enums"
	"jdy/logic/store"
	"jdy/model"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type StoreSalesTotalRes struct {
	Store model.Store `json:"-"` // 店铺

	Name                string          `json:"name"`                  // 店铺名称
	Total               decimal.Decimal `json:"total"`                 // 总业绩
	Silver              decimal.Decimal `json:"silver"`                // 银饰
	Gold                decimal.Decimal `json:"gold"`                  // 足金（件）
	GoldDeduction       decimal.Decimal `json:"gold_deduction"`        // 足金件兑换旧料抵扣
	GoldWeight          decimal.Decimal `json:"gold_weight"`           // 足金克
	GoldWeightDeduction decimal.Decimal `json:"gold_weight_deduction"` // 足金克兑换旧料抵扣
	PieceAccessories    decimal.Decimal `json:"piece_accessories"`     // 计件配件
}

type StoreSalesTotalLogic struct {
	Db  *gorm.DB
	Req *types.StatisticStoreSalesTotalReq
}

func (l *StatisticLogic) StoreSalesTotal(req *types.StatisticStoreSalesTotalReq) (*[]StoreSalesTotalRes, error) {
	var (
		logic = &StoreSalesTotalLogic{
			Db:  model.DB,
			Req: req,
		}
		res []StoreSalesTotalRes
	)

	store_logic := store.StoreLogic{
		Staff: l.Staff,
	}
	stores, err := store_logic.My(&types.StoreListMyReq{})
	if err != nil {
		return nil, err
	}

	for _, store := range *stores {

		def := store.Default(l.Staff.Identity)
		if store.Id == def.Id && store.Name == def.Name {
			continue
		}

		StoreSalesTotalRes := StoreSalesTotalRes{
			Store: store,
			Name:  store.Name,
		}

		products, err := logic.getProducts(&StoreSalesTotalRes)
		if err != nil {
			return nil, err
		}

		for _, product := range products {
			if err := logic.getTotal(&product, &StoreSalesTotalRes); err != nil {
				return nil, err
			}
			if err := logic.getSilver(&product, &StoreSalesTotalRes); err != nil {
				return nil, err
			}
		}

		res = append(res, StoreSalesTotalRes)
	}

	return &res, nil
}

// 获取订单产品列表
func (l *StoreSalesTotalLogic) getProducts(res *StoreSalesTotalRes) ([]model.OrderSalesProduct, error) {
	var (
		products []model.OrderSalesProduct
	)

	db := model.DB.Model(&model.OrderSalesProduct{})
	db = db.Where(&model.OrderSalesProduct{
		Status: enums.OrderSalesStatusComplete,
	})
	db = db.Scopes(model.DurationCondition(l.Req.Duration))
	db = db.Where("store_id = ?", res.Store.Id)
	db = model.OrderSalesProduct{}.Preloads(db)

	if err := db.Find(&products).Error; err != nil {
		return nil, errors.New("获取总业绩失败")
	}

	return products, nil
}

// 获取总业绩
func (l *StoreSalesTotalLogic) getTotal(product *model.OrderSalesProduct, res *StoreSalesTotalRes) error {
	switch product.Type {
	case enums.ProductTypeFinished:
		res.Total = res.Total.Add(product.Finished.Price.Round(2))
	case enums.ProductTypeOld:
		res.Total = res.Total.Add(product.Old.RecyclePrice.Round(2))
	case enums.ProductTypeAccessorie:
		res.Total = res.Total.Add(product.Accessorie.Price.Round(2))
	}

	return nil
}

// 获取银饰数量
func (l *StoreSalesTotalLogic) getSilver(product *model.OrderSalesProduct, res *StoreSalesTotalRes) error {
	switch product.Type {
	case enums.ProductTypeFinished:
		if product.Finished.Product.Class == enums.ProductClassFinishedSilver {
			res.Silver = res.Silver.Add(product.Finished.Price.Round(2))
		}
	case enums.ProductTypeOld:
		if product.Old.Product.Class == enums.ProductClassOldSilver {
			res.Silver = res.Silver.Add(product.Old.RecyclePrice.Round(2))
		}
	}

	return nil
}
