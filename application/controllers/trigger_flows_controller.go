package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"github.com/gin-gonic/gin"
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

	if err := tfc.serverCore.ConfigStore.TriggerFlowRepo.InsertTriggerFlow(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
