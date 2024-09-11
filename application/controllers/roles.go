package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/common"
	"ifttt/manager/domain/roles"
	"ifttt/manager/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type roleController struct {
	serverCore *core.ServerCore
}

func newRoleController(serverCore *core.ServerCore) *roleController {
	return &roleController{
		serverCore: serverCore,
	}
}

func (r *roleController) GetAllPermissions(c *gin.Context) {
	common.ResponseHandler(c, common.ResponseConfig{Data: *r.serverCore.Routes})
}

func (r *roleController) UpdateUserRoles(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, roles.UpdateUserRolesRequest{})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*roles.UpdateUserRolesRequest)

	if user, err := r.serverCore.ConfigStore.UserRepo.GetUser(reqBody.Email, user.DecodeUser); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if user == nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["UserNotFound"]})
		return
	}

	if _, err := r.serverCore.ConfigStore.CasbinEnforcer.DeleteRolesForUser(reqBody.Email); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	if _, err := r.serverCore.ConfigStore.CasbinEnforcer.AddRolesForUser(reqBody.Email, reqBody.Roles); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	if err := r.serverCore.ConfigStore.CasbinEnforcer.SavePolicy(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	if err := r.serverCore.ConfigStore.CasbinEnforcer.LoadPolicy(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (r *roleController) AddUpdateRole(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, roles.UpdateRoleRequest{})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*roles.UpdateRoleRequest)

	policies, err := r.serverCore.ConfigStore.CasbinEnforcer.GetPolicy()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	for _, policy := range policies {
		if policy[0] == reqBody.RoleName {
			_, err := r.serverCore.ConfigStore.CasbinEnforcer.RemovePolicy(common.ConvertStringToInterfaceArray(policy)...)
			if err != nil {
				common.HandleErrorResponse(c, err)
				return
			}
		}
	}

	for _, perm := range reqBody.Permissions {
		permString := CreatePermission(perm.Path, perm.Method)
		if _, found := lo.Find(*r.serverCore.Permissions, func(p string) bool {
			return p == permString
		}); !found {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["PermissionNotFound"]})
			return
		}

		if _, err := r.serverCore.ConfigStore.CasbinEnforcer.AddPolicy(
			reqBody.RoleName, perm.Path, perm.Method,
		); err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
	}

	if err := r.serverCore.ConfigStore.CasbinEnforcer.SavePolicy(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	r.serverCore.ConfigStore.CasbinEnforcer.ClearPolicy()
	if err := r.serverCore.ConfigStore.CasbinEnforcer.LoadPolicy(); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}
