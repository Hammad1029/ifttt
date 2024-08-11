package controllers

import (
	"generic/application/middlewares"
	"generic/application/schemas"
	"generic/domain/api"
	infrastructure "generic/infrastructure/init"
	"generic/utils"

	"github.com/gin-gonic/gin"
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
	_, found, err := a.store.ApiRepository.GetApiByGroupAndName(reqBody.ApiGroup, reqBody.ApiName)
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	} else if found {
		utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["ApiAlreadyExists"]})
		return
	}

	// create api struct
	api := api.Api{
		ApiGroup:       reqBody.ApiGroup,
		ApiName:        reqBody.ApiName,
		ApiPath:        reqBody.ApiPath,
		ApiDescription: reqBody.ApiDescription,
		ApiRequest:     reqBody.ApiRequest,
		StartRules:     reqBody.StartRules,
		Rules:          reqBody.Rules,
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
