package constant

const (
	TenderCompanyStatus_Eligible             string = ""
	TenderCompanyStatus_Waiting              string = "waiting_for_quotation"
	TenderCompanyStatus_CCSubmitted          string = "cc_submitted"
	TenderCompanyStatus_QuotationSubmitted   string = "quotation_submitted"
	TenderCompanyStatus_QuotationReSubmitted string = "quotation_resubmitted"
	TenderCompanyStatus_InReview             string = "in_review"
	TenderCompanyStatus_NegoSession          string = "negotiation_session"
	TenderCompanyStatus_Missed               string = "missed"
	TenderCompanyStatus_Cancelled            string = "cancelled"
	TenderCompanyStatus_Won                  string = "won"
	TenderCompanyStatus_Withdrawn            string = "withdrawn"
	TenderCompanyStatus_Lose                 string = "lose"
)

var TenderCompanyStatusMap = map[string]string{
	TenderCompanyStatus_Waiting:              "Waiting for Quotation",
	TenderCompanyStatus_CCSubmitted:          "CC Submitted",
	TenderCompanyStatus_QuotationSubmitted:   "Quotation Submitted",
	TenderCompanyStatus_QuotationReSubmitted: "Quotation Resubmitted",
	TenderCompanyStatus_InReview:             "In Review",
	TenderCompanyStatus_NegoSession:          "Negotiation Session",
	TenderCompanyStatus_Missed:               "You Missed the Tender",
	TenderCompanyStatus_Cancelled:            "Tender Cancelled",
	TenderCompanyStatus_Won:                  "You Won the Tender!",
	TenderCompanyStatus_Withdrawn:            "Tender Withdrawn",
	TenderCompanyStatus_Lose:                 "You've not been Selected",
}

var TenderCompanyProgressStatuses = []string{
	TenderCompanyStatus_CCSubmitted,
	TenderCompanyStatus_QuotationSubmitted,
	TenderCompanyStatus_QuotationReSubmitted,
	TenderCompanyStatus_InReview,
	TenderCompanyStatus_NegoSession,
}

var TenderCompanyHistoryStatuses = []string{
	TenderCompanyStatus_Missed,
	TenderCompanyStatus_Cancelled,
	TenderCompanyStatus_Won,
	TenderCompanyStatus_Withdrawn,
	TenderCompanyStatus_Lose,
}

var TenderHasBidder = []string{
	TenderCompanyStatus_CCSubmitted,
	TenderCompanyStatus_QuotationSubmitted,
	TenderCompanyStatus_QuotationReSubmitted,
	TenderCompanyStatus_InReview,
	TenderCompanyStatus_NegoSession,
}

var EligibleNotifCancel = []string{
	TenderCompanyStatus_CCSubmitted,
	TenderCompanyStatus_QuotationSubmitted,
	TenderCompanyStatus_QuotationReSubmitted,
	TenderCompanyStatus_InReview,
	TenderCompanyStatus_NegoSession,
}
