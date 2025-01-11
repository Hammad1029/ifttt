package responseprofiles

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (p *Profile) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.MappedCode, validation.Required, validation.Length(2, 3), is.Digit),
		validation.Field(&p.HttpStatus, validation.Required),
		validation.Field(&p.Internal),
		validation.Field(&p.Code, validation.Required, validation.By(
			func(value any) error {
				v := value.(Field)
				return v.Validate()
			})),
		validation.Field(&p.Description, validation.Required, validation.By(
			func(value any) error {
				v := value.(Field)
				return v.Validate()
			})),
		validation.Field(&p.Data, validation.Required, validation.By(
			func(value any) error {
				v := value.(Field)
				return v.Validate()
			})),
	)
}

func (p *Field) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Key, is.Alphanumeric),
		validation.Field(&p.Default),
		validation.Field(&p.Disabled),
	)
}
