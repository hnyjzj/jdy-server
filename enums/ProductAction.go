package enums

import "errors"

/* 产品操作 */
// 入库、更新、报损、旧料转成品、报损转成品、报损转旧料、直接调出、门店间调拨、调拨撤销、调拨确认、开单、撒销开单、退货
type ProductAction int

const (
	ProductActionEntry           ProductAction = iota + 1 // 入库
	ProductActionUpdate                                   // 更新
	ProductActionDamage                                   // 报损
	ProductActionOldToNew                                 // 旧料转成品
	ProductActionDamageToNew                              // 报损转成品
	ProductActionDamageToOld                              // 报损转旧料
	ProductActionDirectOut                                // 直接调出
	ProductActionTransfer                                 // 门店间调拨
	ProductActionTransferCancel                           // 调拨撤销
	ProductActionTransferConfirm                          // 调拨确认
	ProductActionOrder                                    // 开单
	ProductActionOrderCancel                              // 撒销开单
	ProductActionReturn                                   // 退货
	ProductActionEntryCancel                              // 入库撤销
)

var ProductActionMap = map[ProductAction]string{
	ProductActionEntry:           "入库",
	ProductActionUpdate:          "更新",
	ProductActionDamage:          "报损",
	ProductActionOldToNew:        "旧料转成品",
	ProductActionDamageToNew:     "报损转成品",
	ProductActionDamageToOld:     "报损转旧料",
	ProductActionDirectOut:       "直接调出",
	ProductActionTransfer:        "门店间调拨",
	ProductActionTransferCancel:  "调拨撤销",
	ProductActionTransferConfirm: "调拨确认",
	ProductActionOrder:           "开单",
	ProductActionOrderCancel:     "撒销开单",
	ProductActionReturn:          "退货",
}

func (p ProductAction) ToMap() any {
	return ProductActionMap
}

func (p ProductAction) InMap() error {
	if _, ok := ProductActionMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
