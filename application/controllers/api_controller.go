package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/api"
	requestvalidator "ifttt/manager/domain/request_validator"
	"ifttt/manager/domain/resolvable"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type apiController struct {
	serverCore *core.ServerCore
}

func newApiController(serverCore *core.ServerCore) *apiController {
	return &apiController{
		serverCore: serverCore,
	}
}

func (ac *apiController) GetAll(c *gin.Context) {
	apis, err := ac.serverCore.ConfigStore.APIRepo.GetAllApis()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: apis})
}

func (ac *apiController) GetDetails(c *gin.Context) {
	var reqBody api.GetDetailsRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	api, err := ac.serverCore.ConfigStore.APIRepo.GetApiByNameOrPath(reqBody.Name, reqBody.Path)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: api})
}

func (ac *apiController) Create(c *gin.Context) {
	var reqBody api.CreateApiRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if api, err := ac.serverCore.ConfigStore.APIRepo.GetApiByNameOrPath(reqBody.Name, reqBody.Path); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if api != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["APIAlreadyExists"]})
		return
	}

	tNames := lo.Map(reqBody.Triggers, func(t triggerflow.TriggerConditionRequest, _ int) string {
		return t.Trigger
	})

	requiredTFlows, err := ac.serverCore.ConfigStore.TriggerFlowRepo.GetTriggerFlowsByNames(tNames)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if len(*requiredTFlows) != len(reqBody.Triggers) {
		common.ResponseHandler(c,
			common.ResponseConfig{Response: common.Responses["TriggerFlowNotFound"]})
		return
	}

	if err := requestvalidator.GenerateAll(&reqBody.Request); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	if manipulated, err := resolvable.ManipulateMap(reqBody.PreConfig, ac.serverCore.ResolvableDependencies); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		reqBody.PreConfig = manipulated
	}

	for idx, tc := range reqBody.Triggers {
		if err := tc.Manipulate(ac.serverCore.ResolvableDependencies); err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
		reqBody.Triggers[idx] = tc
	}

	if err := ac.serverCore.ConfigStore.APIRepo.InsertApi(&reqBody, requiredTFlows); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
