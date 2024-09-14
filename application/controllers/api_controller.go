package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/api"

	"github.com/gin-gonic/gin"
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
	if ok := validateAndBind(c, &reqBody); ok {
		return
	}

	api, err := ac.serverCore.ConfigStore.APIRepo.GetApiDetailsByNameAndPath(reqBody.Name, reqBody.Path)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: api})
}

func (ac *apiController) Create(c *gin.Context) {
	var reqBody api.CreateApiRequest
	if ok := validateAndBind(c, &reqBody); ok {
		return
	}

	if api, err := ac.serverCore.ConfigStore.APIRepo.GetApiByNameOrPath(reqBody.Name, reqBody.Path); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if api != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["APIAlreadyExists"]})
		return
	}

	if err := ac.serverCore.ConfigStore.APIRepo.InsertApi(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
