package controllers

import "ifttt/manager/application/core"

type AllController struct {
	ApiController    *apiController
	TablesController *schemaController
}

func NewAllController(serverCore *core.ServerCore) *AllController {
	return &AllController{
		ApiController:    newApiController(serverCore),
		TablesController: newSchemaController(serverCore),
	}
}
