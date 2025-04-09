package enums

// MemberIntegralRuleType 积分规则类型
type MemberIntegralRuleType int

const (
	MemberIntegralRuleTypeFinished MemberIntegralRuleType = iota + 1 //成品
	MemberIntegralRuleTypeOld                                        // 旧料
)
