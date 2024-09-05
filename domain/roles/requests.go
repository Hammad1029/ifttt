package roles

import "ifttt/manager/common"

type UpdateUserRolesRequest struct {
	Email string   `mapstructure:"email" json:"email"`
	Roles []string `mapstructure:"roles" json:"roles"`
}

type UpdateRoleRequest struct {
	RoleName    string                   `mapstructure:"roleName" json:"roleName"`
	Permissions []common.RouteDefinition `mapstructure:"permission" json:"permission"`
}
