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

	tIds := lo.Map(reqBody.TriggerFlows, func(t triggerflow.TriggerConditionRequest, _ int) uint {
		return t.Trigger
	})

	if requiredTFlows, err := cC.serverCore.ConfigStore.TriggerFlowRepo.GetTriggerFlowsByIds(tIds); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if len(*requiredTFlows) != len(reqBody.TriggerFlows) {
		common.ResponseHandler(c,
			common.ResponseConfig{Response: common.Responses["TriggerFlowNotFound"]})
		return
	}

	if err := resolvable.ManipulateArray(
		lo.MapToSlice(reqBody.PreConfig,
			func(_ string, r resolvable.Resolvable) resolvable.Resolvable {
				return r
			})); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	if err := cC.serverCore.ConfigStore.CronRepo.InsertCron(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
