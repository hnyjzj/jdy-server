package enums

import "errors"

/* 供应商 */
// 金美福、周大生、老庙、潮宏基、金大福、周六福
type ProductSupplier int

const (
	ProductSupplierJMF ProductSupplier = iota + 1 // 金美福
	ProductSupplierZDS                            // 周大生
	ProductSupplierLMG                            // 老庙
	ProductSupplierCHJ                            // 潮宏基
	ProductSupplierJDF                            // 金大福
	ProductSupplierZLF                            // 周六福
)

var ProductSupplierMap = map[ProductSupplier]string{
	ProductSupplierJMF: "金美福",
	ProductSupplierZDS: "周大生",
	ProductSupplierLMG: "老庙",
	ProductSupplierCHJ: "潮宏基",
	ProductSupplierJDF: "金大福",
	ProductSupplierZLF: "周六福",
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
