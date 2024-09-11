package common

import "github.com/gin-gonic/gin"

type RouteDefinition struct {
	Path          string            `mapstructure:"path" json:"path"`
	Method        string            `mapstructure:"method" json:"method"`
	Description   string            `mapstructure:"description" json:"description"`
	Children      []RouteDefinition `mapstructure:"group" json:"group"`
	Authenticated bool              `mapstructure:"authenticated" json:"authenticated"`
	Authorized    bool              `mapstructure:"authorized" json:"authorized"`
	HandlerFunc   gin.HandlerFunc
}
