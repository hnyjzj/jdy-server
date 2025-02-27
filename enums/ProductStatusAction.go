package enums

import "errors"

/* 产品状态 */
// 全部、入库、更新、报损、旧料转成品、报损转成品、报损转旧料、直接调出、门店间调拨、调拨撤销、调拨确认、开单、撒销开单、退货
type ProductStatusAction int

const (
	ProductStatusActionAll             ProductStatusAction = iota // 全部
	ProductStatusActionEntry                                      // 入库
	ProductStatusActionUpdate                                     // 更新
	ProductStatusActionDamage                                     // 报损
	ProductStatusActionOldToNew                                   // 旧料转成品
	ProductStatusActionDamageToNew                                // 报损转成品
	ProductStatusActionDamageToOld                                // 报损转旧料
	ProductStatusActionDirectOut                                  // 直接调出
	ProductStatusActionTransfer                                   // 门店间调拨
	ProductStatusActionTransferCancel                             // 调拨撤销
	ProductStatusActionTransferConfirm                            // 调拨确认
	ProductStatusActionOrder                                      // 开单
	ProductStatusActionOrderCancel                                // 撒销开单
	ProductStatusActionReturn                                     // 退货
)

var ProductStatusActionMap = map[ProductStatusAction]string{
	ProductStatusActionAll:             "全部",
	ProductStatusActionEntry:           "入库",
	ProductStatusActionUpdate:          "更新",
	ProductStatusActionDamage:          "报损",
	ProductStatusActionOldToNew:        "旧料转成品",
	ProductStatusActionDamageToNew:     "报损转成品",
	ProductStatusActionDamageToOld:     "报损转旧料",
	ProductStatusActionDirectOut:       "直接调出",
	ProductStatusActionTransfer:        "门店间调拨",
	ProductStatusActionTransferCancel:  "调拨撤销",
	ProductStatusActionTransferConfirm: "调拨确认",
	ProductStatusActionOrder:           "开单",
	ProductStatusActionOrderCancel:     "撒销开单",
	ProductStatusActionReturn:          "退货",
}

func (p ProductStatusAction) ToMap() any {
	return ProductStatusActionMap
}

func (p ProductStatusAction) InMap() error {
	if _, ok := ProductStatusActionMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
