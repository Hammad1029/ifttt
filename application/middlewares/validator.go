package middlewares

import (
	"errors"
	"generic/utils"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type validatorInterface interface {
	Validate() error
}

func Validator(c *gin.Context, schema any) (error, any) {
	if schema != nil {
		schemaType := reflect.TypeOf(schema)
		schemaInstance := reflect.New(schemaType).Interface()
		if err := c.ShouldBindJSON(&schemaInstance); err == nil {
			if validator, ok := schemaInstance.(validatorInterface); ok {
				if err = validator.Validate(); err != nil {
					errorRespond(c, err)
					return errors.New("validation failed"), nil
				}
			}
		} else {
			errorRespond(c, err)
			return errors.New("couldn't bind JSON"), nil
		}
		return nil, schemaInstance
	}
	return nil, nil
}

func NestedValidator(target any, fieldRules ...*validation.FieldRules) *validation.FieldRules {
	return validation.Field(target, validation.By(func(value any) error {
		valueV := reflect.Indirect(reflect.ValueOf(value))
		if valueV.CanAddr() {
			addr := valueV.Addr().Interface()
			return validation.ValidateStruct(addr, fieldRules...)
		}
		return validation.ValidateStruct(target, fieldRules...)
	}))
}

func errorRespond(c *gin.Context, err error) {
	ValidationErrors := make(map[string]any)
	ValidationErrors["Validation Errors"] = strings.Split(err.Error(), "; ")
	utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["BadRequest"], Data: ValidationErrors})
}
