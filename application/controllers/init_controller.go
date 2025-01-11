package controllers

import (
	"ifttt/manager/application/core"
)

type AllController struct {
	ApiController              *apiController
	TriggerFlowsController     *triggerFlowsController
	RulesController            *rulesController
	SchemaController           *ormSchemaController
	AuthController             *authController
	UserController             *userController
	RoleController             *roleController
	CronController             *cronController
	ResponseProfilesController *responseProfilesController
}

func NewAllController(serverCore *core.ServerCore) *AllController {
	return &AllController{
		ApiController:              newApiController(serverCore),
		TriggerFlowsController:     newTriggerFlowsController(serverCore),
		RulesController:            newRulesController(serverCore),
		SchemaController:           newOrmSchemaController(serverCore),
		AuthController:             newAuthController(serverCore),
		UserController:             newUserController(serverCore),
		RoleController:             newRoleController(serverCore),
		CronController:             newCronController(serverCore),
		ResponseProfilesController: newResponseProfilesController(serverCore),
	}
}
