package roles

type UpdateUserRolesRequest struct {
	Email string   `mapstructure:"email" json:"email"`
	Roles []string `mapstructure:"roles" json:"roles"`
}

type UpdateRoleRequest struct {
	RoleName    string              `mapstructure:"roleName" json:"roleName"`
	Permissions []permissionRequest `mapstructure:"permissions" json:"permissions"`
}

type permissionRequest struct {
	Path   string `mapstructure:"path" json:"path"`
	Method string `mapstructure:"method" json:"method"`
}
