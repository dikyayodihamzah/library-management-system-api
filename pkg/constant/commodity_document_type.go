package constant

const (
	DocumentType_Legal        string = "legal"
	DocumentType_Support      string = "support"
	DocumentType_HSE          string = "hse"
	DocumentType_Company      string = "company"
	DocumentType_Registration string = "registration"
)

func DocumentType() []string {
	return []string{
		DocumentType_Legal,
		DocumentType_Support,
		DocumentType_HSE,
		DocumentType_Company,
	}
}
