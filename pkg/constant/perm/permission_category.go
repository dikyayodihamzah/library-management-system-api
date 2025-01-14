package perm

const (
	Category_Permission = "permission"
	Category_Vendor     = "vendor"
	Category_Commodity  = "commodity"
	Category_Quotation  = "quotation"
	Category_Tender     = "tender"
)

var CategoryMap = map[string]string{
	Category_Permission: "Permission Management",
	Category_Vendor:     "Vendor Management",
	Category_Commodity:  "Commodity Management",
	Category_Quotation:  "Quotation Management",
	Category_Tender:     "Tender Management",
}
