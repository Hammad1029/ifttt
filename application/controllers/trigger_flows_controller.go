package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type triggerFlowsController struct {
	serverCore *core.ServerCore
}

func newTriggerFlowsController(serverCore *core.ServerCore) *triggerFlowsController {
	return &triggerFlowsController{
		serverCore: serverCore,
	}
}

func (tfc *triggerFlowsController) GetAll(c *gin.Context) {
	tFlows, err := tfc.serverCore.ConfigStore.TriggerFlowRepo.GetAllTriggerFlows()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: tFlows})
}

func (tfc *triggerFlowsController) GetDetails(c *gin.Context) {
	var reqBody triggerflow.GetDetailsRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	tFlow, err := tfc.serverCore.ConfigStore.TriggerFlowRepo.GetTriggerFlowDetailsByName(reqBody.Name)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: tFlow})
}

func (tfc *triggerFlowsController) Create(c *gin.Context) {
	var reqBody triggerflow.CreateTriggerFlowRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if existing, err := tfc.serverCore.ConfigStore.TriggerFlowRepo.GetTriggerFlowByName(reqBody.Name); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if existing != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["FlowAlreadyExists"]})
		return
	}

	if requiredRules, err := tfc.serverCore.ConfigStore.RuleRepo.GetRulesByIds(reqBody.Rules); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if len(*requiredRules) != len(reqBody.Rules) {
		common.ResponseHandler(c,
			common.ResponseConfig{Response: common.Responses["TriggerFlowRulesNotFound"]})
		return
	}

	for _, flow := range reqBody.BranchFlows {
		if _, found := lo.Find(reqBody.Rules, func(r uint) bool {
			return r == flow.Rule
		}); !found {
			common.ResponseHandler(c,
				common.ResponseConfig{Response: common.Responses["InvalidBranchFlow"]})
			return
		}

		if ok := lo.Every(common.RuleAllowedReturns, lo.Keys(flow.States)); !ok {
			common.ResponseHandler(c,
				common.ResponseConfig{Response: common.Responses["InvalidBranchFlow"]})
			return
		}
	}

	if _, ok := reqBody.BranchFlows[reqBody.StartState]; !ok {
		common.ResponseHandler(c,
			common.ResponseConfig{Response: common.Responses["InvalidBranchFlow"]})
		return
	}

	if err := tfc.serverCore.ConfigStore.TriggerFlowRepo.InsertTriggerFlow(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
