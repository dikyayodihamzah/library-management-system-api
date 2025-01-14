package constant

const (
	RiskLevel_Low    int = 1
	RiskLevel_Medium int = 2
	RiskLevel_High   int = 3
)

func RiskLevel() []int {
	return []int{
		RiskLevel_Low,
		RiskLevel_Medium,
		RiskLevel_High,
	}
}
