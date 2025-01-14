package lib

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// VALIDATOR validate request body
var VALIDATOR *validator.Validate = validator.New()

// BodyParser with validation
func BodyParser(c *fiber.Ctx, payload interface{}) error {
	if err := c.BodyParser(payload); nil != err {
		return err
	}

	return VALIDATOR.Struct(payload)
}

// func QueryParser(c *fiber.Ctx, out interface{}) error {
// 	if err := c.QueryParser(out); err != nil {
// 		return err
// 	}

// 	val := reflect.ValueOf(out)

// 	if val.Kind() == reflect.Pointer {
// 		val = reflect.Indirect(val)
// 	}

// 	for i := 0; i < val.Type().NumField(); i++ {
// 		t := val.Type().Field(i)
// 		fieldName := t.Name

// 		if t.Type.Kind() == reflect.Slice {
// 			if value := val.Field(i).Slice() {

// 			}
// 		}
// 		if t.Type.Kind() == reflect.Struct {
// 			keys = append(keys, PrintFields(val.Field(i).Interface(), k)...)

// 		}
// 		switch jsonTag := t.Tag.Get("json"); jsonTag {
// 		case "-", "", "deleted_at":
// 		default:
// 			parts := strings.Split(jsonTag, ",")
// 			name := parts[0]
// 			if name == "" {
// 				name = fieldName
// 			}
// 			keys = append(keys, name)
// 		}

// 	}

// 	return nil
// }
