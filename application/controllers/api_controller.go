package controllers

import (
	"ifttt/manager/application/middlewares"
	"ifttt/manager/application/schemas"
	"ifttt/manager/domain/api"
	infrastructure "ifttt/manager/infrastructure/init"
	"ifttt/manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type apiController struct {
	store *infrastructure.DbStore
}

func newApiController(store *infrastructure.DbStore) *apiController {
	return &apiController{
		store: store,
	}
}

func (a *apiController) CreateApi(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, schemas.CreateApiRequest{})
	if err != nil {
		return
	}
	reqBody := reqBodyAny.(*schemas.CreateApiRequest)

	// check if api of this name already exists
	_, found, err := a.store.ApiRepository.GetApiByGroupAndName(reqBody.Group, reqBody.Name)
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
	if err := a.store.ApiRepository.InsertApi(apiSerialized); err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	utils.ResponseHandler(c, utils.ResponseConfig{})
}

func (a *apiController) GetApis(c *gin.Context) {
	apis, err := a.store.ApiRepository.GetAllApis()
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	utils.ResponseHandler(c, utils.ResponseConfig{Data: apis})
}
