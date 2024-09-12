package roles

type UpdateUserRolesRequest struct {
	Email string   `mapstructure:"email" json:"email"`
	Roles []string `mapstructure:"roles" json:"roles"`
}

type UpdateRoleRequest struct {
	RoleName    string              `mapstructure:"roleName" json:"roleName"`
	Permissions []PermissionVerbose `mapstructure:"permissions" json:"permissions"`
}
