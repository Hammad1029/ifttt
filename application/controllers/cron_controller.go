package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/cron"
	"ifttt/manager/domain/resolvable"
	triggerflow "ifttt/manager/domain/trigger_flow"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type cronController struct {
	serverCore *core.ServerCore
}

func newCronController(serverCore *core.ServerCore) *cronController {
	return &cronController{
		serverCore: serverCore,
	}
}

func (cC *cronController) GetAll(c *gin.Context) {
	cronJobs, err := cC.serverCore.ConfigStore.CronRepo.GetAllCrons()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: cronJobs})
}

func (cC *cronController) GetByName(c *gin.Context) {
	var reqBody cron.GetCronRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	cronJob, err := cC.serverCore.ConfigStore.CronRepo.GetCronByName(reqBody.Name)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: cronJob})
}

func (cC *cronController) Create(c *gin.Context) {
	var reqBody cron.CreateCronRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if cron, err := cC.serverCore.ConfigStore.CronRepo.GetCronByName(reqBody.Name); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if cron != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["CronAlreadyExists"]})
		return
	}

	tNames := lo.Map(reqBody.Triggers, func(t triggerflow.TriggerConditionRequest, _ int) string {
		return t.Trigger
	})

	requiredTFlows, err := cC.serverCore.ConfigStore.TriggerFlowRepo.GetTriggerFlowsByNames(tNames)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if len(*requiredTFlows) != len(reqBody.Triggers) {
		common.ResponseHandler(c,
			common.ResponseConfig{Response: common.Responses["TriggerFlowNotFound"]})
		return
	}

	if manipulated, err := resolvable.ManipulateMap(reqBody.PreConfig, cC.serverCore.ResolvableDependencies); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		reqBody.PreConfig = manipulated
	}

	for idx, tc := range reqBody.Triggers {
		if err := tc.Manipulate(cC.serverCore.ResolvableDependencies); err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
		reqBody.Triggers[idx] = tc
	}

	if err := cC.serverCore.ConfigStore.CronRepo.InsertCron(&reqBody, requiredTFlows); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
