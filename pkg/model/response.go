package model

type Response struct {
	Code             int         `json:"code"`            // http status
	Message          string      `json:"message"`         // response message
	Status           bool        `json:"status"`          // response status true of false
	Page             *int        `json:"page,omitempty"`  // page number
	Count            *int        `json:"count,omitempty"` // data count
	Total            *int        `json:"total,omitempty"` // total data
	Data             interface{} `json:"data"`
	ErrorDescription *string     `json:"error_description,omitempty" swaggerignore:"true"` // Error description
	Log              *string     `json:"log,omitempty"`
	// ErrorData        *[]ErrorData `json:"error_data,omitempty" swaggerignore:"true"`        // Error fields
}

// ErrorData field error details
type ErrorData struct {
	Name      string      `json:"name" example:"Username"`                                 // data name
	Path      string      `json:"path" example:"user.username"`                            // object property path
	Type      string      `json:"type,omitempty" example:"string"`                         // data type
	Value     interface{} `json:"value,omitempty" swaggertype:"string" example:"jane doe"` // value
	Validator string      `json:"validator" example:"required"`                            // validator type, see [more details](https://github.com/go-playground/validator#baked-in-validations)
	Criteria  interface{} `json:"criteria,omitempty" swaggertype:"number" example:"10"`    // criteria, example: if validator is gte (greater than) and criteria is 10, then it means a maximum of 10
	Message   string      `json:"message" example:"invalid value"`                         // Field message
}

var ValidateMessages = map[string]string{
	"required":    "%v is required %v",
	"gte":         "%v must be greater than or equal %v",
	"gt":          "%v must be greater than %v",
	"lte":         "%v must be less than or equal %v",
	"lt":          "%v must be less than %v",
	"email":       "%v not valid email %v",
	"number":      "%v not valid number %v",
	"oneof":       "%v not valid value : %v",
	"required_if": "%v is required when %v",
}

func (e *Response) Error() string {
	return e.Message
}
