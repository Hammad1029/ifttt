package roles

import (
	"ifttt/manager/common"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type AddUpdateRoleRequest struct {
	RoleName    string              `mapstructure:"roleName" json:"roleName"`
	Permissions []PermissionVerbose `mapstructure:"permissions" json:"permissions"`
	AssignTo    []string            `mapstructure:"assignTo" json:"assignTo"`
}

func (a *AddUpdateRoleRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.RoleName, validation.Required, validation.Length(5, 50), is.Alphanumeric),
		validation.Field(&a.Permissions, validation.Required, validation.Length(1, 0), validation.Each(
			validation.By(func(value interface{}) error {
				perm := value.(PermissionVerbose)
				return validation.ValidateStruct(&perm,
					validation.Field(&perm.Method, validation.Required,
						validation.In(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)),
					validation.Field(&perm.Path, validation.Match(common.RegexEndpoint)),
				)
			}),
		)),
		validation.Field(&a.AssignTo, validation.Each(validation.Required, is.Email)),
	)
}
