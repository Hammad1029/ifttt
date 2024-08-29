package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/domain/api"
	"ifttt/manager/utils"

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
		return
	}
	reqBody := reqBodyAny.(*api.CreateApiRequest)

	// check if api of this name already exists
	_, found, err := a.serverCore.ConfigStore.APIRepo.GetApiByGroupAndName(reqBody.Group, reqBody.Name)
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	} else if found {
		utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ApiAlreadyExists"]})
		return
	}

	// create api struct
	var api api.Api
	if err := mapstructure.Decode(reqBody, &api); err != nil {
		utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ServerError"]})
	}

	// serialize data
	apiSerialized, err := api.Serialize()
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	// insert api
	if err := a.serverCore.ConfigStore.APIRepo.InsertApi(apiSerialized); err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	utils.ResponseHandler(c, utils.ResponseConfig{})
}

func (a *apiController) GetApis(c *gin.Context) {
	apis, err := a.serverCore.ConfigStore.APIRepo.GetAllApis()
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	utils.ResponseHandler(c, utils.ResponseConfig{Data: apis})
}
