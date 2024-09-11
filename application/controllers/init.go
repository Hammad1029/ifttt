package controllers

import (
	"fmt"
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"net/http"
)

type AllController struct {
	ApiController    *apiController
	SchemaController *schemaController
	AuthController   *authController
	UserController   *userController
	RoleController   *roleController
}

func NewAllController(serverCore *core.ServerCore) *AllController {
	return &AllController{
		ApiController:    newApiController(serverCore),
		SchemaController: newSchemaController(serverCore),
		AuthController:   newAuthController(serverCore),
		UserController:   newUserController(serverCore),
		RoleController:   newRoleController(serverCore),
	}
}

func (ac *AllController) GetRoutes() *[]common.RouteDefinition {
	return &[]common.RouteDefinition{
		{
			Path:        "/auth",
			Method:      "GROUP",
			Description: "Authentication Group",
			Children: []common.RouteDefinition{
				{
					Path:          "/login",
					Method:        http.MethodPost,
					Description:   "User Login",
					Authenticated: false,
					Authorized:    false,
					HandlerFunc:   ac.AuthController.Login,
				},
				{
					Path:          "/refresh",
					Method:        http.MethodGet,
					Description:   "Refresh Token",
					Authenticated: false,
					Authorized:    false,
					HandlerFunc:   ac.AuthController.RefreshToken,
				},
				{
					Path:          "/logout",
					Method:        http.MethodGet,
					Description:   "User Logout",
					Authenticated: true,
					Authorized:    false,
					HandlerFunc:   ac.AuthController.Logout,
				},
			},
		},
		{
			Path:        "/apis",
			Method:      "GROUP",
			Description: "APIs Group",
			Children: []common.RouteDefinition{
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create API",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.ApiController.CreateApi,
				},
			},
		},
		{
			Path:        "/schemas",
			Method:      "GROUP",
			Description: "Schemas Group",
			Children: []common.RouteDefinition{
				{
					Path:          "/getAll",
					Method:        http.MethodGet,
					Description:   "Get Schema",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.SchemaController.GetSchema,
				},
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create Table",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.SchemaController.CreateTable,
				},
				{
					Path:          "/update",
					Method:        http.MethodPost,
					Description:   "Update Table",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.SchemaController.UpdateTable,
				},
			},
		},
		{
			Path:        "/users",
			Method:      "GROUP",
			Description: "User Group",
			Children: []common.RouteDefinition{
				{
					Path:          "/get",
					Method:        http.MethodGet,
					Description:   "Get all users",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.UserController.GetAllUsers,
				},
				{
					Path:          "/create",
					Method:        http.MethodPost,
					Description:   "Create User",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.UserController.CreateUser,
				},
			},
		},
		{
			Path:        "/roles",
			Method:      "GROUP",
			Description: "Update Table",
			Children: []common.RouteDefinition{
				{
					Path:          "/getAllPermissions",
					Method:        http.MethodGet,
					Description:   "Get All Permissions",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.RoleController.GetAllPermissions,
				},
				{
					Path:          "/updateUserRoles",
					Method:        http.MethodPost,
					Description:   "Update User Roles",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.RoleController.UpdateUserRoles,
				},
				{
					Path:          "/addUpdateRoles",
					Method:        http.MethodPost,
					Description:   "Add/Update Roles",
					Authenticated: true,
					Authorized:    true,
					HandlerFunc:   ac.RoleController.AddUpdateRole,
				},
			},
		},
	}
}

func CreatePermission(path string, method string) string {
	return fmt.Sprintf("%s:%s", path, method)
}
