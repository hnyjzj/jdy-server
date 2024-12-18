package enums

import "errors"

/* 供应商 */
// 全部、金美福、周大生、老庙、潮宏基
type ProductSupplier int

const (
	ProductSupplierAll ProductSupplier = iota // 全部
	ProductSupplierJMF                        // 金美福
	ProductSupplierZDS                        // 周大生
	ProductSupplierLM                         // 老庙
	ProductSupplierZLF                        // 周六福
	ProductSupplierJDF                        // 金大福
)

var ProductSupplierMap = map[ProductSupplier]string{
	ProductSupplierAll: "全部",
	ProductSupplierJMF: "金美福",
	ProductSupplierZDS: "周大生",
	ProductSupplierLM:  "老庙",
	ProductSupplierZLF: "周六福",
	ProductSupplierJDF: "金大福",
}

func (p ProductSupplier) ToMap() any {
	return ProductSupplierMap
}

func (p ProductSupplier) InMap() error {
	if _, ok := ProductSupplierMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
