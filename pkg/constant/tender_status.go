package constant

const (
	TenderStatus_Preparation = iota + 1
	TenderStatus_Incoming
	TenderStatus_OnLive
	TenderStatus_CloseBid
	TenderStatus_InReview
	TenderStatus_Done
	TenderStatus_Cancelled
	TenderStatus_Recalled
)

func TenderStatus() []int {
	return []int{
		TenderStatus_Preparation,
		TenderStatus_Incoming,
		TenderStatus_OnLive,
		TenderStatus_CloseBid,
		TenderStatus_InReview,
		TenderStatus_Done,
		TenderStatus_Cancelled,
		TenderStatus_Recalled,
	}
}

var MapTenderStatus map[int]string = map[int]string{
	TenderStatus_Preparation: "Preparation Stage",
	TenderStatus_Incoming:    "Incoming Tender",
	TenderStatus_OnLive:      "On Live",
	TenderStatus_CloseBid:    "Close Bid",
	TenderStatus_InReview:    "In Review",
	TenderStatus_Done:        "Done",
	TenderStatus_Cancelled:   "Tender Cancelled",
	TenderStatus_Recalled:    "Recalled Tender",
}
