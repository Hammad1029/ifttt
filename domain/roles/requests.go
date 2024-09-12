package roles

import (
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
		validation.Field(&a.Permissions, validation.Required, validation.Length(5, 50)),
	)
}
