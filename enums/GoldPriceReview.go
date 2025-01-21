package enums

/* 黄金价格状态 */
// 审批、驳回
type GoldPriceReview string

const (
	GoldPriceReviewApproved GoldPriceReview = "gold_price_approval"  // 审批
	GoldPriceReviewRejected GoldPriceReview = "gold_price_rejection" // 驳回
)
