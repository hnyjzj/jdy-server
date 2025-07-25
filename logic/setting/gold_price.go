package setting

import (
	"errors"
	"jdy/enums"
	"jdy/logic"
	"jdy/message"
	"jdy/model"
	"jdy/types"
	"log"

	"gorm.io/gorm"
)

type GoldPriceLogic struct {
	logic.BaseLogic

	IP string
}

// 设置金价
func (l *GoldPriceLogic) Create(req *types.GoldPriceCreateReq) error {
	if len(req.Options) == 0 {
		return errors.New("请至少设置一条金价")
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 添加/更新金价列表
		for _, v := range req.Options {
			var ProductBrand []enums.ProductBrand
			if len(v.ProductBrand) == 0 {
				ProductBrand = enums.ProductBrandJMF.All()
			} else {
				ProductBrand = v.ProductBrand
			}
			// 转换数据结构
			data := model.GoldPrice{
				StoreId:         v.StoreId,
				Price:           v.Price,
				ProductMaterial: v.ProductMaterial,
				ProductType:     v.ProductType,
				ProductBrand:    ProductBrand,
				ProductQuality:  v.ProductQuality,
			}

			if v.Id == "" {
				if err := tx.Create(&data).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&model.GoldPrice{}).Where("id = ?", v.Id).Updates(data).Error; err != nil {
					return err
				}
			}
		}

		// 删除
		for _, v := range req.Deletes {
			if err := tx.Delete(&model.GoldPrice{}, "id = ?", v).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return errors.New("设置金价失败")
	}

	// 发送审批消息
	go func() {
		var store model.Store
		if err := model.DB.Where("id = ?", req.Options[0].StoreId).
			Preload("Staffs").
			First(&store).Error; err != nil {
			log.Printf("获取店铺信息失败: %v\n", err)
			return
		}
		var receiver []string
		for _, v := range store.Staffs {
			receiver = append(receiver, v.Username)
		}
		m := message.NewMessage(l.Ctx)
		m.SendGoldPriceUpdateMessage(&message.GoldPriceMessage{
			ToUser:    receiver,
			StoreName: store.Name,
			StoreId:   store.Id,
			Operator:  l.Staff.Nickname,
		})
	}()

	return nil
}

func (l *GoldPriceLogic) List(req *types.GoldPriceListReq) (*[]model.GoldPrice, error) {
	var (
		gold_price model.GoldPrice
		res        []model.GoldPrice
	)

	db := model.DB.Order("updated_at desc")
	db = gold_price.WhereCondition(db, &types.GoldPriceOptions{StoreId: req.StoreId})
	// 获取列表
	db = db.Order("updated_at desc")
	if err := db.Find(&res).Error; err != nil {
		return nil, errors.New("获取金价列表失败")
	}

	return &res, nil
}
