package eventprofiles

import (
	"ifttt/manager/domain/resolvable"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (p *Profile) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Trigger, validation.Required, validation.Length(2, 3), is.Digit),
		validation.Field(&p.ResponseHTTPStatus, validation.Required),
		validation.Field(&p.UseBody),
		validation.Field(&p.ResponseBody, validation.NotNil, validation.By(
			func(value any) error {
				return resolvable.ValidateIfResolvable(value)
			})),
		validation.Field(&p.Internal),
	)
}
