package server

import (
	"ifttt/manager/application/controllers"
	"ifttt/manager/common"
	"net/http"
)

func getRouteDefinitions(controllers *controllers.AllController) *[]common.RouteDefinition {
	return &[]common.RouteDefinition{
		{
			Path:        "/auth",
			Method:      "GROUP",
			Description: "Authentication",
			Authorized:  false,
			Children: []common.RouteDefinition{
				{
					Path:          "/login",
					Method:        http.MethodPost,
					Description:   "User Login",
					Authenticated: false,
					Authorized:    false,
					HandlerFunc:   controllers.AuthController.Login,
				},
				{
					Path:          "/refresh",
					Method:        http.MethodGet,
					Description:   "Refresh Token",
					Authenticated: false,
					Authorized:    false,
					HandlerFunc:   controllers.AuthController.RefreshToken,
				},
				{
					Path:          "/logout",
					Method:        http.MethodGet,
					Description:   "User Logout",
					Authenticated: true,
					Authorized:    false,
					HandlerFunc:   controllers.AuthController.Logout,
				},
			},
		},
		{
			Path:        "/apis",
			Method:      "GROUP",
			Description: "API Management",
			Authorized:  true,
			Children: []common.RouteDefinition{
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create API",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.ApiController.CreateApi,
				},
			},
		},
		{
			Path:        "/schemas",
			Method:      "GROUP",
			Description: "Schema Management",
			Authorized:  true,
			Children: []common.RouteDefinition{
				{
					Path:          "/getAll",
					Method:        http.MethodGet,
					Description:   "Get Schema",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.SchemaController.GetSchema,
				},
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create Table",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.SchemaController.CreateTable,
				},
				{
					Path:          "/update",
					Method:        http.MethodPost,
					Description:   "Update Table",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.SchemaController.UpdateTable,
				},
			},
		},
		{
			Path:        "/users",
			Method:      "GROUP",
			Description: "User Management",
			Authorized:  true,
			Children: []common.RouteDefinition{
				{
					Path:          "/get",
					Method:        http.MethodGet,
					Description:   "Get all users",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.UserController.GetAllUsers,
				},
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create User",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.UserController.CreateUser,
				},
			},
		},
		{
			Path:        "/roles",
			Method:      "GROUP",
			Description: "Role Management",
			Authorized:  true,
			Children: []common.RouteDefinition{
				{
					Path:          "/getAllRoles",
					Method:        http.MethodGet,
					Description:   "Get All Roles",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.RoleController.GetAllRoles,
				},
				{
					Path:          "/getAllPermissions",
					Method:        http.MethodGet,
					Description:   "Get All Permissions",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.RoleController.GetAllPermissions,
				},
				{
					Path:          "/updateUserRoles",
					Method:        http.MethodPost,
					Description:   "Update User Roles",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.RoleController.UpdateUserRoles,
				},
				{
					Path:          "/addUpdateRoles",
					Method:        http.MethodPost,
					Description:   "Add/Update Roles",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.RoleController.AddUpdateRole,
				},
			},
		},
	}
}
