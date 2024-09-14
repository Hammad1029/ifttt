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
					Path:          "/getAll",
					Method:        http.MethodGet,
					Description:   "Get All APIs",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.ApiController.GetAll,
				},
				{
					Path:          "/getDetails",
					Method:        http.MethodPost,
					Description:   "Get API Details By Name And Path",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.ApiController.GetDetails,
				},
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create API",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.ApiController.Create,
				},
			},
		},
		{
			Path:        "/triggerFlows",
			Method:      "GROUP",
			Description: "Trigger Flow Management",
			Authorized:  true,
			Children: []common.RouteDefinition{
				{
					Path:          "/getAll",
					Method:        http.MethodGet,
					Description:   "Get All Trigger Flows",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.TriggerFlowsController.GetAll,
				},
				{
					Path:          "/getDetails",
					Method:        http.MethodPost,
					Description:   "Get Trigger Flow Details By Name",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.TriggerFlowsController.GetDetails,
				},
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create Trigger Flow",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.TriggerFlowsController.Create,
				},
			},
		},
		{
			Path:        "/rules",
			Method:      "GROUP",
			Description: "Rule Management",
			Authorized:  true,
			Children: []common.RouteDefinition{
				{
					Path:          "/getAll",
					Method:        http.MethodGet,
					Description:   "Get All Rules",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.RulesController.GetAll,
				},
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create Rule",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   controllers.RulesController.Create,
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
