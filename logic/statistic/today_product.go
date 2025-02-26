package statistic

import (
	"database/sql"
	"errors"
	"jdy/enums"
	"jdy/model"
	"jdy/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TodayProductRes struct {
	ProductStockCount int64           `json:"product_stock_count"` // 成品库存件数
	OldStockCount     int64           `json:"old_stock_count"`     // 旧料库存件数
	OldStockWeight    decimal.Decimal `json:"old_stock_weight"`    // 旧料库存金重
	UnsalableCount    int64           `json:"unsalable_count"`     // 滞销货品件数
}

type TodayProductLogic struct {
	Req *types.StatisticTodayProductReq
	Db  *gorm.DB
	Res *TodayProductRes
}

func (l *StatisticLogic) TodayProduct(req *types.StatisticTodayProductReq) (*TodayProductRes, error) {
	logic := &TodayProductLogic{
		Req: req,
		Res: &TodayProductRes{},
		Db:  model.DB,
	}

	if err := logic.getProductStockCount(); err != nil {
		return nil, err
	}

	if err := logic.getOldStock(); err != nil {
		return nil, err
	}

	if err := logic.getUnsalableCount(); err != nil {
		return nil, err
	}

	return logic.Res, nil
}

// 获取成品库存件数
func (l *TodayProductLogic) getProductStockCount() error {
	var (
		count sql.NullInt64
	)

	if err := l.Db.Model(&model.Product{}).
		Where(&model.Product{
			Status:  enums.ProductStatusNormal,
			Type:    enums.ProductTypeFinished,
			StoreId: l.Req.StoreId,
		}).Select("sum(stock) as count").Scan(&count).Error; err != nil {

		return errors.New("获取成品库存件数失败")
	}

	if count.Valid {
		l.Res.ProductStockCount = count.Int64
	} else {
		l.Res.ProductStockCount = 0
	}

	return nil
}

type oldStock struct {
	Count  sql.NullInt64   `json:"count"`  // 件数
	Weight sql.NullFloat64 `json:"weight"` // 金重
}

// 获取旧料库存件数
func (l *TodayProductLogic) getOldStock() error {
	var (
		res oldStock
	)

	if err := l.Db.Model(&model.Product{}).
		Where(&model.Product{
			Status:  enums.ProductStatusNormal,
			StoreId: l.Req.StoreId,
			Type:    enums.ProductTypeOld,
		}).Select("sum(stock) as count, sum(weight_metal) as weight").Scan(&res).Error; err != nil {

		return errors.New("获取旧料库存件数失败")
	}

	if res.Count.Valid {
		l.Res.OldStockCount = res.Count.Int64
	} else {
		l.Res.OldStockCount = 0
	}

	if res.Weight.Valid {
		l.Res.OldStockWeight = decimal.NewFromFloat(res.Weight.Float64)
	} else {
		l.Res.OldStockWeight = decimal.Zero
	}

	return nil
}

// 获取滞销货品件数
func (l *TodayProductLogic) getUnsalableCount() error {
	var (
		count sql.NullInt64
	)

	if err := l.Db.Model(&model.Product{}).
		Where(&model.Product{
			Status:  enums.ProductStatusNormal,
			StoreId: l.Req.StoreId,
		}).
		Where("type in (?)", []enums.ProductType{enums.ProductTypeFinished, enums.ProductTypeOld}). // 成品、旧料
		// 创建时间大于 6 个月，即为滞销货品
		Where("created_at < DATE_SUB(NOW(), INTERVAL 6 MONTH)").
		Select("sum(stock) as count").Scan(&count).Error; err != nil {

		return errors.New("获取滞销货品件数失败")
	}

	if count.Valid {
		l.Res.UnsalableCount = count.Int64
	} else {
		l.Res.UnsalableCount = 0
	}

	return nil
}
