package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/resolvable"
	"ifttt/manager/domain/rule"

	"github.com/gin-gonic/gin"
)

type rulesController struct {
	serverCore *core.ServerCore
}

func newRulesController(serverCore *core.ServerCore) *rulesController {
	return &rulesController{
		serverCore: serverCore,
	}
}

func (rc *rulesController) GetAll(c *gin.Context) {
	var reqBody rule.GetRulesRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	rules, err := rc.serverCore.ConfigStore.RuleRepo.GetRulesLikeName(reqBody.Name)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: rules})
}

func (rc *rulesController) Create(c *gin.Context) {
	var reqBody rule.CreateRuleRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if existing, err := rc.serverCore.ConfigStore.RuleRepo.GetRuleByName(reqBody.Name); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if existing != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["RuleAlreadyExists"]})
		return
	}

	if reqBody.Switch.Default.Return != common.RuleDefaultReturn {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["InvalidReturnValue"]})
		return
	}

	if manipulated, err := resolvable.ManipulateArray(&reqBody.Pre, rc.serverCore.ResolvableDependencies); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		reqBody.Pre = *manipulated
	}

	if err := reqBody.Switch.Manipulate(rc.serverCore.ResolvableDependencies); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	if manipulated, err := resolvable.ManipulateArray(&reqBody.Finally, rc.serverCore.ResolvableDependencies); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		reqBody.Finally = *manipulated
	}

	if err := rc.serverCore.ConfigStore.RuleRepo.InsertRule(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
