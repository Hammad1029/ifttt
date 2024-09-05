package common

import "net/http"

type RouteDefinition struct {
	Path        string `mapstructure:"path" json:"path"`
	Method      string `mapstructure:"method" json:"method"`
	Description string `mapstructure:"description" json:"description"`
}

type routeIdx int

const (
	AuthGroup routeIdx = iota + 1
	AuthLogin
	AuthRefresh
	AuthLogout

	ApisGroup
	ApisGet
	ApisCreate

	SchemasGroup
	SchemasGet
	SchemasCreateTable
	SchemasUpdateTable

	UserGroup
	UserCreate

	RolesGroup
	RolesGetAllPermissions
	RolesUpdateUserRoles
	RolesAddUpdateRole
)

var RouteDefinitions = map[routeIdx]RouteDefinition{
	AuthGroup:   {Path: "/auth", Method: "GROUP", Description: "Authentication Group"},
	AuthLogin:   {Path: "/login", Method: http.MethodPost, Description: "User Login"},
	AuthRefresh: {Path: "/refresh", Method: http.MethodPost, Description: "Refresh Token"},
	AuthLogout:  {Path: "/logout", Method: http.MethodGet, Description: "User Logout"},

	ApisGroup:  {Path: "/apis", Method: "GROUP", Description: "APIs Group"},
	ApisGet:    {Path: "/", Method: http.MethodGet, Description: "Get APIs"},
	ApisCreate: {Path: "/create", Method: http.MethodPost, Description: "Create API"},

	SchemasGroup:       {Path: "/schemas", Method: "GROUP", Description: "Schemas Group"},
	SchemasGet:         {Path: "/", Method: http.MethodGet, Description: "Get Schema"},
	SchemasCreateTable: {Path: "/create", Method: http.MethodPost, Description: "Create Table"},
	SchemasUpdateTable: {Path: "/update", Method: http.MethodPost, Description: "Update Table"},

	UserGroup:  {Path: "/users", Method: "GROUP", Description: "User Group"},
	UserCreate: {Path: "/create", Method: http.MethodPost, Description: "Create User"},

	RolesGroup:             {Path: "/roles", Method: "GROUP", Description: "Update Table"},
	RolesGetAllPermissions: {Path: "/getAllPermissions", Method: http.MethodGet, Description: "Get All Permissions"},
	RolesUpdateUserRoles:   {Path: "/updateUserRoles", Method: http.MethodPost, Description: "Update User Roles"},
	RolesAddUpdateRole:     {Path: "/addUpdateRoles", Method: http.MethodPost, Description: "Add/Update Roles"},
}

func (r *RouteDefinition) GetRouteByPathAndMethod() *RouteDefinition {
	for _, route := range RouteDefinitions {
		if route.Path == r.Path && route.Method == r.Method {
			return &route
		}
	}
	return nil
}

func GetRoutePath(idx routeIdx) string {
	return RouteDefinitions[idx].Path
}
