package constant

const (
	AccountStatus_Waiting    int = 1
	AccountStatus_Verified   int = 2
	AccountStatus_Unverified int = 3
	AccountStatus_Suspended  int = 4
	AccountStatus_ReVerified int = 5
)

func AccountStatus() []int {
	return []int{
		AccountStatus_Waiting,
		AccountStatus_Verified,
		AccountStatus_Unverified,
		AccountStatus_Suspended,
		AccountStatus_ReVerified,
	}
}

var EligibleRegister []int = []int{
	AccountStatus_Verified,
	AccountStatus_ReVerified,
}
