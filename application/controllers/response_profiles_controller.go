package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	responseprofiles "ifttt/manager/domain/response_profiles"

	"github.com/gin-gonic/gin"
)

type responseProfilesController struct {
	serverCore *core.ServerCore
}

func newResponseProfilesController(serverCore *core.ServerCore) *responseProfilesController {
	return &responseProfilesController{
		serverCore: serverCore,
	}
}

func (r *responseProfilesController) AddProfile(c *gin.Context) {
	var reqBody responseprofiles.Profile
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if existing, err := r.serverCore.ConfigStore.ResponseProfileRepo.
		GetProfilesByInternalAndCode(true, reqBody.MappedCode); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if reqBody.Internal && existing != nil && len(*existing) > 0 {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
		return
	} else if !reqBody.Internal && (existing == nil || len(*existing) == 0) {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["NotFound"]})
		return
	} else if reqBody.Internal {
		if err := r.serverCore.ConfigStore.ResponseProfileRepo.AddProfile(&reqBody, 0); err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
	} else if !reqBody.Internal {
		if existing, err := r.serverCore.ConfigStore.ResponseProfileRepo.
			GetProfilesByInternalAndCode(false, reqBody.MappedCode); err != nil {
			common.HandleErrorResponse(c, err)
			return
		} else if existing != nil && len(*existing) > 0 {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
			return
		}
		if err := r.serverCore.ConfigStore.ResponseProfileRepo.AddProfile(&reqBody, (*existing)[0].ID); err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (r *responseProfilesController) GetAllProfiles(c *gin.Context) {
	if profiles, err := r.serverCore.ConfigStore.ResponseProfileRepo.GetAllInternalProfiles(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		common.ResponseHandler(c, common.ResponseConfig{Data: profiles})
		return
	}
}
