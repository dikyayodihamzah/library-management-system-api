package constant

const (
	NotificationPrefix string = "notification-"

	// preparation
	NotifType_NotifyFeedback string = "notify_feedback"
	NotifType_Hm1ClosePrep   string = "day_min_1_close_preparation"

	// incoming
	NotifType_Hm1Onlive     string = "day_min_1_auto_onlive"
	NotifType_OverdueOnLive string = "overdue_on_live"

	// on live
	NotifType_Hm1CloseBid       string = "day_min_1_close_bidding"
	NotifType_Hm1CloseBidVendor string = "day_min_1_close_bidding_vendor"

	// review
	NotifType_TenderReview      string = "notify_tender_review"
	NotifType_StartNegotiation  string = "start_negotiation"
	NotifType_EndNegotiation    string = "end_negotiation"
	NotifType_Hm1EndNegotiation string = "day_min_1_end_negotiation"
	NotifType_Hm1CloseReview    string = "day_min_1_close_review"

	// approval
	NotifType_OverdueApproval     string = "overdue_approval"
	NotifType_NotifyApproval      string = "notify_approval"
	NotifType_Hm1CloseApproval    string = "day_min_1_close_approval"
	NotifType_UnfinishedApproval  string = "unfinished_approval"
	NotifType_UnfinishedApproval7 string = "unfinished_approval_7"

	// edit tender
	NotifType_TenderChangeVisibility string = "tender_change_visibility"
	NotifType_CommodityRegistration  string = "commodity_registration"

	// vendor
	NotifType_Verify          string = "notify_verify"
	NotifType_Verify7         string = "notify_verify_7"
	NotifType_Verify14        string = "notify_verify_14"
	NotifType_Verify30        string = "notify_verify_30"
	NotifType_ResetPassword   string = "notify_reset_password"
	NotifType_ResetPassword7  string = "notify_reset_password_7"
	NotifType_ResetPassword14 string = "notify_reset_password_14"
	NotifType_ResetPassword30 string = "notify_reset_password_30"
)
