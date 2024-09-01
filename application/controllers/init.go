package controllers

import "ifttt/manager/application/core"

type allController struct {
	ApiController    *apiController
	TablesController *schemaController
	AuthController   *authController
}

func NewAllController(serverCore *core.ServerCore) *allController {
	return &allController{
		ApiController:    newApiController(serverCore),
		TablesController: newSchemaController(serverCore),
		AuthController:   newAuthController(serverCore),
	}
}
