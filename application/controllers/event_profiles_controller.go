package controllers

import (
	"fmt"
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	eventprofiles "ifttt/manager/domain/event_profiles"
	"ifttt/manager/domain/resolvable"

	"github.com/gin-gonic/gin"
)

type eventProfilesController struct {
	serverCore *core.ServerCore
}

func newEventProfilesController(serverCore *core.ServerCore) *eventProfilesController {
	return &eventProfilesController{
		serverCore: serverCore,
	}
}

func (r *eventProfilesController) AddProfile(c *gin.Context) {
	var reqBody eventprofiles.Profile
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	existing, err := r.serverCore.ConfigStore.EventProfileRepo.
		GetProfilesByInternalAndTrigger(true, reqBody.Trigger)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if reqBody.Internal && existing != nil && len(*existing) > 0 {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
		return
	} else if !reqBody.Internal && (existing == nil || len(*existing) == 0) {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["NotFound"]})
		return
	}

	if manipulated, err := resolvable.ManipulateIfResolvable(
		reqBody.ResponseBody, r.serverCore.ResolvableDependencies); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		common.HandleErrorResponse(c, fmt.Errorf("could not cast manipulated to map[string]any"))
		return
	} else {
		reqBody.ResponseBody = mapped
	}

	if reqBody.Internal {
		if err := r.serverCore.ConfigStore.EventProfileRepo.AddProfile(&reqBody, 0); err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
	} else if !reqBody.Internal {
		if existing, err := r.serverCore.ConfigStore.EventProfileRepo.
			GetProfilesByInternalAndTrigger(false, reqBody.Trigger); err != nil {
			common.HandleErrorResponse(c, err)
			return
		} else if existing != nil && len(*existing) > 0 {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
			return
		}
		if err := r.serverCore.ConfigStore.EventProfileRepo.AddProfile(&reqBody, (*existing)[0].ID); err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (r *eventProfilesController) GetAllProfiles(c *gin.Context) {
	if profiles, err := r.serverCore.ConfigStore.EventProfileRepo.GetAllInternalProfiles(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		common.ResponseHandler(c, common.ResponseConfig{Data: profiles})
		return
	}
}
