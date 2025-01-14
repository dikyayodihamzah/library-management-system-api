package constant

const (
	ReviewProgress_Review    = "review"    // tender in review, approval not started (all approver in waiting status)
	ReviewProgress_Approval  = "approval"  // tender in approval, started when first approver is need_recommendation until all approver give approval
	ReviewProgress_Completed = "completed" // tender in approval, when all approver already given an approval (no more waiting or need recommendation)
)
