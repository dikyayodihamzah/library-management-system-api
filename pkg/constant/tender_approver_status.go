package constant

const (
	ApproverStatus_NeedFeedback       string = "need_feedback"
	ApproverStatus_Waiting            string = "waiting"
	ApproverStatus_NeedRecommendation string = "need_recommendation"
	ApproverStatus_RecommendationSend string = "recommendation_send"
	ApproverStatus_Disapprove         string = "disapprove"
	ApproverStatus_Recall             string = "recall"
	ApproverStatus_Cancelled          string = "cancelled"
)

var ActionRecommend []string = []string{
	ApproverStatus_Disapprove,
	ApproverStatus_RecommendationSend,
}

// var MapApproverStatus map[string]string = map[string]string{
// 	"waiting":             "Waiting",
// 	"need_recommendation": "Need Recommendation",
// 	"recommendation_send": "Recommendation Send",
// 	"disapprove":          "Disapproved",
// 	"recall":              "Tender Recalled",
// 	"cancelled":           "Tender Cancelled",
// }
