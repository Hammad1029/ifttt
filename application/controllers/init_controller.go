package controllers

import (
	"ifttt/manager/application/core"
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
