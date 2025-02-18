package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/cron"
	"net/http"

	"github.com/gin-gonic/gin"
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
	var reqBody cron.Cron
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

	if api, err := cC.serverCore.ConfigStore.APIRepo.GetApiByNameOrPath(reqBody.ApiName, ""); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if api == nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["NotFound"]})
		return
	} else if api.Method != http.MethodGet {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Response{Code: "21", Description: "only get allowed"}})
		return
	} else if err := cC.serverCore.ConfigStore.CronRepo.InsertCron(&reqBody, api.ID); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
