package constant

const (
	QuotationFieldType_FreeText = iota + 1
	QuotationFieldType_CurrencyIDR
	QuotationFieldType_Number
	QuotationFieldType_Text
	QuotationFieldType_Document
	QuotationFieldType_CurrencyUSD
)

var QuotationFieldTypeMap = map[int]string{
	QuotationFieldType_FreeText:    "Free text",
	QuotationFieldType_CurrencyIDR: "Currency (IDR)",
	QuotationFieldType_Number:      "Number",
	QuotationFieldType_Text:        "Text",
	QuotationFieldType_Document:    "Document",
	QuotationFieldType_CurrencyUSD: "Currency (USD)",
}
