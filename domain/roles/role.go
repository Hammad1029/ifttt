package roles

import "fmt"

type Role struct {
	RoleName    string              `mapstructure:"roleName" json:"roleName"`
	Permissions []PermissionVerbose `mapstructure:"permissions" json:"permissions"`
}

type PermissionVerbose struct {
	Path   string `mapstructure:"path" json:"path"`
	Method string `mapstructure:"method" json:"method"`
}

func (p *PermissionVerbose) CreatePermission() string {
	return fmt.Sprintf("%s:%s", p.Path, p.Method)
}
