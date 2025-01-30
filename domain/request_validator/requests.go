package requestvalidator

import (
	"fmt"
	"ifttt/manager/common"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mitchellh/mapstructure"
)

func (r *RequestParameter) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.DataType, validation.Required,
			validation.In(dataTypeText, dataTypeNumber, dataTypeBoolean, dataTypeArray, dataTypeMap)),
		validation.Field(&r.Required),
		validation.Field(&r.InternalTag),
		validation.Field(&r.Config, validation.Required, validation.By(
			func(value interface{}) error {
				var validator common.Validatable
				switch r.DataType {
				case dataTypeText:
					validator = &textValue{}
				case dataTypeNumber:
					validator = &numberValue{}
				case dataTypeBoolean:
					validator = &booleanValue{}
				case dataTypeArray:
					validator = &arrayValue{}
				case dataTypeMap:
					validator = &mapValue{}
				default:
					return validation.NewError("datatype_not_found", fmt.Sprintf("datatype %s not found", r.DataType))
				}
				data := value.(map[string]any)
				if err := mapstructure.Decode(data, &validator); err != nil {
					return validation.NewInternalError(err)
				}
				if err := validator.Validate(); err != nil {
					return err
				}
				return nil
			})),
	)
}

func (t *textValue) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Alpha),
		validation.Field(&t.Numeric),
		validation.Field(&t.Special),
		validation.Field(&t.Minimum),
		validation.Field(&t.Maximum),
		validation.Field(&t.In),
	)
}

func (n *numberValue) Validate() error {
	return validation.ValidateStruct(n,
		validation.Field(&n.Minimum),
		validation.Field(&n.Maximum),
		validation.Field(&n.In),
	)
}

func (b *booleanValue) Validate() error {
	return nil
}

func (a *arrayValue) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Minimum),
		validation.Field(&a.Maximum),
		validation.Field(&a.OfType, validation.Required, validation.By(
			func(value interface{}) error {
				if casted, ok := value.(*RequestParameter); !ok {
					return validation.NewError("validator-not-casted",
						fmt.Sprintf("could not cast validator for %s", dataTypeArray))
				} else {
					return casted.Validate()
				}
			})),
	)
}

func (m *mapValue) Validate() error {
	for _, v := range *m {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}
