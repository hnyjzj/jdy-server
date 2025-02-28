package enums

import "errors"

/* 产品操作 */
// 全部、入库、更新、报损、旧料转成品、报损转成品、报损转旧料、直接调出、门店间调拨、调拨撤销、调拨确认、开单、撒销开单、退货
type ProductAction int

const (
	ProductActionAll             ProductAction = iota // 全部
	ProductActionEntry                                // 入库
	ProductActionUpdate                               // 更新
	ProductActionDamage                               // 报损
	ProductActionOldToNew                             // 旧料转成品
	ProductActionDamageToNew                          // 报损转成品
	ProductActionDamageToOld                          // 报损转旧料
	ProductActionDirectOut                            // 直接调出
	ProductActionTransfer                             // 门店间调拨
	ProductActionTransferCancel                       // 调拨撤销
	ProductActionTransferConfirm                      // 调拨确认
	ProductActionOrder                                // 开单
	ProductActionOrderCancel                          // 撒销开单
	ProductActionReturn                               // 退货
)

var ProductActionMap = map[ProductAction]string{
	ProductActionAll:             "全部",
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
