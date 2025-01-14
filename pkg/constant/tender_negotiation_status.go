package constant

const (
	NegoStatus_InNego  = "in_nego"
	NegoStatus_Waiting = "waiting"
	NegoStatus_Missed  = "missed"
	NegoStatus_Done    = "done"
)

func NegoStatus() []string {
	return []string{
		NegoStatus_InNego,
		NegoStatus_Waiting,
		NegoStatus_Missed,
		NegoStatus_Done,
	}
}

var NegoStatusMap map[string]string = map[string]string{
	NegoStatus_InNego:  "In Negotiation",
	NegoStatus_Waiting: "Waiting Negotiation Approval",
	NegoStatus_Missed:  "Missed Negotiation",
	NegoStatus_Done:    "Done Negotiation",
}
