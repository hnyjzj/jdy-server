package sync

import (
	"jdy/enums"
	"jdy/model"

	"gorm.io/gorm"
)

type PaymentLogic struct {
	SyncLogic
}

func (l *PaymentLogic) Payments() error {
	var payments []model.OrderPayment
	db := model.DB.Model(&model.OrderPayment{})
	db = db.Where(&model.OrderPayment{
		Status: false,
	})
	db = model.OrderPayment{}.Preloads(db)
	if err := db.Find(&payments).Error; err != nil {
		return err
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		for _, payment := range payments {
			switch payment.OrderType {
			case enums.OrderTypeSales:
				{
					if payment.OrderSales.Status != enums.OrderSalesStatusWaitPay &&
						payment.OrderSales.Status != enums.OrderSalesStatusCancel {
						if err := tx.Model(&model.OrderPayment{}).Where("id = ?", payment.Id).Update("status", true).Error; err != nil {
							return err
						}
					}
				}
			case enums.OrderTypeDeposit:
				{
					if payment.OrderDeposit.Status != enums.OrderDepositStatusWaitPay &&
						payment.OrderDeposit.Status != enums.OrderDepositStatusCancel {
						if err := tx.Model(&model.OrderPayment{}).Where("id = ?", payment.Id).Update("status", true).Error; err != nil {
							return err
						}
					}
				}
			case enums.OrderTypeRepair:
				{
					if payment.OrderRepair.Status != enums.OrderRepairStatusWaitPay &&
						payment.OrderRepair.Status != enums.OrderRepairStatusCancel {
						if err := tx.Model(&model.OrderPayment{}).Where("id = ?", payment.Id).Update("status", true).Error; err != nil {
							return err
						}
					}
				}
			}
		}

		return nil

	}); err != nil {
		return err
	}

	return nil
}
