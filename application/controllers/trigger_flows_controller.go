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

	if requiredRules, err := tfc.serverCore.ConfigStore.RuleRepo.GetRulesByIds(reqBody.StartRules); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if len(*requiredRules) != len(reqBody.StartRules) {
		common.ResponseHandler(c,
			common.ResponseConfig{Response: common.Responses["TriggerFlowStartRulesNotFound"]})
		return
	}

	if requiredRules, err := tfc.serverCore.ConfigStore.RuleRepo.GetRulesByIds(reqBody.AllRules); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if len(*requiredRules) != len(reqBody.AllRules) {
		common.ResponseHandler(c,
			common.ResponseConfig{Response: common.Responses["TriggerFlowAllRulesNotFound"]})
		return
	}

	if len(lo.Intersect(reqBody.AllRules, lo.Keys(reqBody.BranchFlows))) != len(reqBody.AllRules) {
		common.ResponseHandler(c,
			common.ResponseConfig{Response: common.Responses["InvalidBranchFlow"]})
		return
	}

	jumpIds := lo.FlatMap(lo.Values(reqBody.BranchFlows), func(bF []triggerflow.BranchFlow, _ int) []uint {
		return lo.Map(bF, func(f triggerflow.BranchFlow, _ int) uint {
			return f.Jump
		})
	})
	if !lo.Every(reqBody.AllRules, jumpIds) {
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
