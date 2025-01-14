package constant

const (
	TenderVisibility_Internal = "internal"
	TenderVisibility_Private  = "private"
	TenderVisibility_Public   = "public"

	TenderSourceTab_Available = "available"
	TenderSourceTab_Progress  = "progress"
	TenderSourceTab_History   = "history"
)

func TenderVisibility() []string {
	return []string{
		TenderVisibility_Internal,
		TenderVisibility_Private,
		TenderVisibility_Public,
	}
}

var MapTenderVisibility map[string]string = map[string]string{
	TenderVisibility_Internal: "Internal",
	TenderVisibility_Private:  "Private",
	TenderVisibility_Public:   "Public",
}

const (
	TenderPublicType_Vendor    = "vendor"
	TenderPublicType_Commodity = "commodity"
)
