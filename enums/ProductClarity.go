package enums

import "errors"

/* 产品净度 */
// IF、VVS、VVS1、VVS2、VS、VS1、VS2、SI、SI1、SI2、SI3、P、FL、LC
type ProductClarity int

const (
	ProductClarityIF   ProductClarity = iota + 1 // IF
	ProductClarityVVS                            // VVS
	ProductClarityVVS1                           // VVS1
	ProductClarityVVS2                           // VVS2
	ProductClarityVS                             // VS
	ProductClarityVS1                            // VS1
	ProductClarityVS2                            // VS2
	ProductClaritySI                             // SI
	ProductClaritySI1                            // SI1
	ProductClaritySI2                            // SI2
	ProductClaritySI3                            // SI3
	ProductClarityP                              // P
	ProductClarityFL                             // FL
	ProductClarityLC                             // LC
)

var ProductClarityMap = map[ProductClarity]string{
	ProductClarityIF:   "IF",
	ProductClarityVVS:  "VVS",
	ProductClarityVVS1: "VVS1",
	ProductClarityVVS2: "VVS2",
	ProductClarityVS:   "VS",
	ProductClarityVS1:  "VS1",
	ProductClarityVS2:  "VS2",
	ProductClaritySI:   "SI",
	ProductClaritySI1:  "SI1",
	ProductClaritySI2:  "SI2",
	ProductClaritySI3:  "SI3",
	ProductClarityP:    "P",
	ProductClarityFL:   "FL",
	ProductClarityLC:   "LC",
}

func (p ProductClarity) ToMap() any {
	return ProductClarityMap
}

func (p ProductClarity) InMap() error {
	if _, ok := ProductClarityMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
