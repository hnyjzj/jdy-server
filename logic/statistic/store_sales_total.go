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

func (*StatisticLogic) StoreSalesTotal(req *types.StatisticStoreSalesTotalReq) (*[]StoreSalesTotalRes, error) {
	var (
		stores []model.Store
		logic  = &StoreSalesTotalLogic{
			Db:  model.DB,
			Req: req,
		}
		res []StoreSalesTotalRes
	)

	sdb := logic.Db.Model(&model.Store{})
	if err := sdb.Find(&stores).Error; err != nil {
		return nil, err
	}

	for _, store := range stores {
		StoreSalesTotalRes := StoreSalesTotalRes{
			Store: store,
			Name:  store.Name,
		}

		if err := logic.getTotal(&StoreSalesTotalRes); err != nil {
			return nil, err
		}
		if err := logic.getSilver(&StoreSalesTotalRes); err != nil {
			return nil, err
		}
		if err := logic.getGold(&StoreSalesTotalRes); err != nil {
			return nil, err
		}
		// if err := logic.getGoldDeduction(&StoreSalesTotalRes); err != nil {
		// 	return nil, err
		// }
		if err := logic.getGoldWeight(&StoreSalesTotalRes); err != nil {
			return nil, err
		}
		// if err := logic.getGoldWeightDeduction(&StoreSalesTotalRes); err != nil {
		// 	return nil, err
		// }
		if err := logic.getPieceAccessories(&StoreSalesTotalRes); err != nil {
			return nil, err
		}

		res = append(res, StoreSalesTotalRes)
	}

	return &res, nil
}

func (l *StoreSalesTotalLogic) getTotal(res *StoreSalesTotalRes) error {
	var (
		db    = model.DB
		total sql.NullFloat64
	)

	if err := db.Model(&model.Order{}).
		Scopes(model.DurationCondition(l.Req.Duration)).
		Where(&model.Order{
			StoreId: res.Store.Id,
			Status:  enums.OrderStatusComplete,
		}).Select("sum(amount_pay) as total").Scan(&total).Error; err != nil {
		return errors.New("获取总业绩失败")
	}

	if total.Valid {
		res.Total = decimal.NewFromFloat(total.Float64)
	} else {
		res.Total = decimal.NewFromFloat(0)
	}

	return nil
}

func (l *StoreSalesTotalLogic) getWhereDb() *gorm.DB {
	db := model.DB.Model(&model.OrderProduct{})
	db = db.
		Joins("JOIN products ON order_products.product_id = products.id").
		Where("order_products.status = ?", enums.OrderStatusComplete).
		Scopes(model.DurationCondition(l.Req.Duration, "order_products.created_at"))

	return db
}

func (l *StoreSalesTotalLogic) getSilver(res *StoreSalesTotalRes) error {
	var (
		silver sql.NullFloat64
	)

	if err := l.getWhereDb().
		Where("products.class =?", enums.ProductClassSilver).
		Select("SUM(order_products.amount)").
		Scan(&silver).Error; err != nil {
		return errors.New("获取银饰数量失败")
	}

	if silver.Valid {
		res.Silver = decimal.NewFromFloat(silver.Float64)
	} else {
		res.Silver = decimal.NewFromFloat(0)
	}

	return nil
}

func (l *StoreSalesTotalLogic) getGold(res *StoreSalesTotalRes) error {
	var (
		gold sql.NullFloat64
	)

	if err := l.getWhereDb().
		Where("products.class =?", enums.ProductClassGoldPiece).
		Select("SUM(order_products.amount)").
		Scan(&gold).Error; err != nil {
		return errors.New("获取足金数量失败")
	}

	if gold.Valid {
		res.Gold = decimal.NewFromFloat(gold.Float64)
	} else {
		res.Gold = decimal.NewFromFloat(0)
	}

	return nil
}
func (l *StoreSalesTotalLogic) getGoldWeight(res *StoreSalesTotalRes) error {
	var (
		goldWeight sql.NullFloat64
	)

	if err := l.getWhereDb().
		Where("products.class =?", enums.ProductClassGoldKg).
		Select("SUM(order_products.amount)").
		Scan(&goldWeight).Error; err != nil {
		return errors.New("获取金重数量失败")
	}

	if goldWeight.Valid {
		res.GoldWeight = decimal.NewFromFloat(goldWeight.Float64)
	} else {
		res.GoldWeight = decimal.NewFromFloat(0)
	}

	return nil
}

func (l *StoreSalesTotalLogic) getPieceAccessories(res *StoreSalesTotalRes) error {
	var (
		pieceAccessories sql.NullFloat64
	)

	if err := l.getWhereDb().
		Where("products.type =?", enums.ProductTypeAccessories).
		// Where("products.class =?", enums.ProductClassPieceAccessories).
		Select("SUM(order_products.amount)").
		Scan(&pieceAccessories).Error; err != nil {
		return errors.New("获取配件数量失败")
	}

	if pieceAccessories.Valid {
		res.PieceAccessories = decimal.NewFromFloat(pieceAccessories.Float64)
	} else {
		res.PieceAccessories = decimal.NewFromFloat(0)
	}

	return nil
}
