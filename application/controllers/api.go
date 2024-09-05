package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/common"
	"ifttt/manager/domain/api"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type apiController struct {
	serverCore *core.ServerCore
}

func newApiController(serverCore *core.ServerCore) *apiController {
	return &apiController{
		serverCore: serverCore,
	}
}

func (a *apiController) CreateApi(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, api.CreateApiRequest{})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*api.CreateApiRequest)

	// check if api of this name already exists
	foundApis, err := a.serverCore.ConfigStore.APIRepo.GetApisByGroupAndName(reqBody.Group, reqBody.Name)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if foundApis != nil && len(*foundApis) > 0 {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["ApiAlreadyExists"]})
		return
	}

	// create api struct
	var api api.Api
	if err := mapstructure.Decode(reqBody, &api); err != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["ServerError"]})
	}

	// insert api
	if err := a.serverCore.ConfigStore.APIRepo.InsertApi(&api); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (a *apiController) GetApis(c *gin.Context) {
	apis, err := a.serverCore.ConfigStore.APIRepo.GetAllApis()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: apis})
}
