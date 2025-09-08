package today

import (
	"database/sql"
	"errors"
	"jdy/enums"
	"jdy/model"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductReq struct {
	DataReq
	StoreId string `json:"store_id" binding:"required"` // 门店ID
}

type ProductRes struct {
	ProductStockCount int64           `json:"product_stock_count"` // 成品库存件数
	OldStockCount     int64           `json:"old_stock_count"`     // 旧料库存件数
	OldStockWeight    decimal.Decimal `json:"old_stock_weight"`    // 旧料库存金重
	UnsalableCount    int64           `json:"unsalable_count"`     // 滞销货品件数
}

type ProductLogic struct {
	Req *ProductReq
	Db  *gorm.DB
	Res *ProductRes

	endtime time.Time
}

func (l *ToDayLogic) Product(req *ProductReq) (*ProductRes, error) {
	logic := &ProductLogic{
		Req: req,
		Res: &ProductRes{},
		Db:  model.DB,
	}

	_, endtime, err := req.Duration.GetTime(time.Now(), req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	logic.endtime = endtime

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
func (l *ProductLogic) getProductStockCount() error {
	var (
		count int64
	)

	db := l.Db.Model(&model.ProductFinished{})
	db = db.Where(&model.ProductFinished{
		StoreId: l.Req.StoreId,
	})
	db = db.Where("status IN (?)", []enums.ProductStatus{
		enums.ProductStatusNormal,
		enums.ProductStatusAllocate,
	})
	db = db.Where("enter_time <= ?", l.endtime)

	if err := db.Count(&count).Error; err != nil {
		return errors.New("获取成品库存件数失败")
	}

	l.Res.ProductStockCount = count

	return nil
}

type oldStock struct {
	Count  sql.NullInt64   `json:"count"`  // 件数
	Weight sql.NullFloat64 `json:"weight"` // 金重
}

// 获取旧料库存件数
func (l *ProductLogic) getOldStock() error {
	var (
		res oldStock
	)

	db := l.Db.Model(&model.ProductOld{})
	db = db.Where(&model.ProductOld{
		Status:  enums.ProductStatusNormal,
		StoreId: l.Req.StoreId,
	})
	db = db.Where("created_at <= ?", l.endtime)
	db = db.Select("COUNT(id) as count, SUM(weight_metal) as weight")

	if err := db.Scan(&res).Error; err != nil {
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
func (l *ProductLogic) getUnsalableCount() error {
	var product_count sql.NullInt64

	db := l.Db.Model(&model.ProductFinished{})
	db = db.Where(&model.ProductFinished{
		StoreId: l.Req.StoreId,
	})
	db = db.Where("status IN (?)", []enums.ProductStatus{
		enums.ProductStatusNormal,
		enums.ProductStatusAllocate,
	})
	db = db.Select("COUNT(id) as count")

	// 滞销货品：enter_time 对比 req.end_time 超过180天
	end_time, err := time.ParseInLocation(time.RFC3339, l.Req.EndTime, time.Now().Location())
	if err != nil {
		return errors.New("时间格式错误")
	}
	db = db.Where("DATEDIFF(?, enter_time) > 180", end_time)

	if err := db.Scan(&product_count).Error; err != nil {
		return errors.New("获取滞销货品件数失败")
	}

	if product_count.Valid {
		l.Res.UnsalableCount = product_count.Int64
	} else {
		l.Res.UnsalableCount = 0
	}

	return nil
}
