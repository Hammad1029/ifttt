package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/configuration"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type configurationController struct {
	serverCore *core.ServerCore
}

func newConfigurationController(serverCore *core.ServerCore) *configurationController {
	return &configurationController{
		serverCore: serverCore,
	}
}

func (r *configurationController) AddResponseProfile(c *gin.Context) {
	var reqBody configuration.ResponseProfile
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if existing, err := r.serverCore.ConfigStore.ResponseProfileRepo.GetProfilesByName(reqBody.Name); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if existing != nil && len(*existing) > 0 {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
		return
	}

	if err := configuration.ValidateMapWithInternalTags(
		&reqBody.BodyFormat, &r.serverCore.ConfigStore.InternalTagRepo,
	); err != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["NotFound"]})
		return
	}

	if err := r.serverCore.ConfigStore.ResponseProfileRepo.AddProfile(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (r *configurationController) GetAllResponseProfiles(c *gin.Context) {
	if profiles, err := r.serverCore.ConfigStore.ResponseProfileRepo.GetAllProfiles(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		common.ResponseHandler(c, common.ResponseConfig{Data: profiles})
		return
	}
}

func (r *configurationController) AddInternalTag(c *gin.Context) {
	var reqBody configuration.InternalTagRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if pTag, err := r.serverCore.ConfigStore.InternalTagRepo.GetByIDOrName(0, reqBody.Name); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if pTag != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
		return
	}

	groups, err := r.serverCore.ConfigStore.InternalTagRepo.GetAllGroups()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	filteredGroups := lo.Filter(*groups, func(g configuration.InternalTagGroup, _ int) bool {
		return lo.ContainsBy(reqBody.Groups, func(id uint) bool {
			return g.ID == id
		})
	})
	if len(filteredGroups) != len(reqBody.Groups) {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["NotFound"]})
		return
	}

	addTag := configuration.InternalTag{
		Name:     reqBody.Name,
		Groups:   filteredGroups,
		All:      reqBody.All,
		Reserved: reqBody.Reserved,
	}
	if err := r.serverCore.ConfigStore.InternalTagRepo.Add(&addTag); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	common.ResponseHandler(c, common.ResponseConfig{})
}

func (r *configurationController) GetAllInternalTags(c *gin.Context) {
	pTags, err := r.serverCore.ConfigStore.InternalTagRepo.GetAll()
	if err != nil {
		common.HandleErrorResponse(c, err)
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: pTags})
}

func (r *configurationController) GetAllInternalTagGroups(c *gin.Context) {
	groups, err := r.serverCore.ConfigStore.InternalTagRepo.GetAllGroups()
	if err != nil {
		common.HandleErrorResponse(c, err)
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: groups})
}

func (r *configurationController) AddInternalTagGroup(c *gin.Context) {
	var reqBody configuration.InternalTagGroup
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if group, err := r.serverCore.ConfigStore.InternalTagRepo.GetGroupByName(reqBody.Name); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if group != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["AlreadyExists"]})
		return
	}

	if err := r.serverCore.ConfigStore.InternalTagRepo.AddGroup(&reqBody); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	common.ResponseHandler(c, common.ResponseConfig{})
}
