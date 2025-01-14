package constant

const (
	VendorStatus_Unregistered         int = 1
	VendorStatus_WaitingRegistration  int = 2
	VendorStatus_NeedApproval         int = 3
	VendorStatus_WaitingRevision      int = 4
	VendorStatus_ApprovalCiculation   int = 5
	VendorStatus_RevisionUpdatesLegal int = 6
	VendorStatus_RevisionUpdatesHSE   int = 7
	VendorStatus_Registered           int = 8
	VendorStatus_Rejected             int = 9
)

func VendorStatus() []int {
	return []int{
		VendorStatus_Unregistered,
		VendorStatus_WaitingRegistration,
		VendorStatus_NeedApproval,
		VendorStatus_WaitingRevision,
		VendorStatus_ApprovalCiculation,
		VendorStatus_RevisionUpdatesLegal,
		VendorStatus_RevisionUpdatesHSE,
		VendorStatus_Registered,
		VendorStatus_Rejected,
	}
}
