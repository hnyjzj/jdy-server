package enums

import (
	"errors"
)

/* 会员等级 */
// 全部、银卡、金卡、钻石卡
type MemberLevel int

const (
	MemberLevelAll     MemberLevel = iota // 全部
	MemberLevelNone                       // 无
	MemberLevelSilver                     // 银卡
	MemberLevelGold                       // 金卡
	MemberLevelDiamond                    // 钻石卡
)

var MemberLevelMap = map[MemberLevel]string{
	MemberLevelAll:     "全部",
	MemberLevelNone:    "无",
	MemberLevelSilver:  "银卡",
	MemberLevelGold:    "金卡",
	MemberLevelDiamond: "钻石卡",
}

func (p MemberLevel) ToMap() any {
	return MemberLevelMap
}

func (p MemberLevel) InMap() error {
	if _, ok := MemberLevelMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

// func (s MemberLevel) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		Value int    `json:"value"`
// 		Desc  string `json:"desc"`
// 	}{
// 		Value: int(s),
// 		Desc:  MemberLevelMap[s],
// 	})
// }
