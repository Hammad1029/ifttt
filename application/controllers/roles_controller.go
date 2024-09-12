package controllers

import (
	"ifttt/manager/application/core"
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
	filteredPermissions := lo.Filter(*r.serverCore.Routes, func(r common.RouteDefinition, _ int) bool {
		return r.Authorized
	})

	common.ResponseHandler(c, common.ResponseConfig{Data: filteredPermissions})
}

func (r *roleController) AddUpdateRole(c *gin.Context) {
	var reqBody *roles.AddUpdateRoleRequest
	if ok := validateRequest(c, &reqBody); !ok {
		return
	}

	var permissions [][]string
	for _, perm := range reqBody.Permissions {
		permString := perm.CreatePermission()
		if _, found := lo.Find(*r.serverCore.Permissions, func(p string) bool {
			return p == permString
		}); !found {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["PermissionNotFound"]})
			return
		}
		permissions = append(permissions, []string{perm.Path, perm.Method})
	}

	if _, err := r.serverCore.ConfigStore.CasbinEnforcer.AddPermissionsForUser(
		reqBody.RoleName, permissions...,
	); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	if users, err := r.serverCore.ConfigStore.CasbinEnforcer.GetUsersForRole(reqBody.RoleName); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else {
		for _, u := range users {
			if _, err := r.serverCore.ConfigStore.CasbinEnforcer.DeleteRoleForUser(u, reqBody.RoleName); err != nil {
				common.HandleErrorResponse(c, err)
				return
			}
		}
	}

	for _, u := range reqBody.AssignTo {
		user, err := r.serverCore.ConfigStore.UserRepo.GetUser(u, user.DecodeUser)
		if err != nil {
			common.HandleErrorResponse(c, err)
			return
		} else if user == nil {
			common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["UserNotFound"]})
			return
		}
		if _, err := r.serverCore.ConfigStore.CasbinEnforcer.AddRoleForUser(u, reqBody.RoleName); err != nil {
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

func (rc *roleController) GetAllRoles(c *gin.Context) {
	var domainRoles []roles.Role

	allRoles, err := rc.serverCore.ConfigStore.CasbinEnforcer.GetAllRoles()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	for _, r := range allRoles {
		newRole := roles.Role{RoleName: r}
		rolePermissions, err := rc.serverCore.ConfigStore.CasbinEnforcer.GetPermissionsForUser(r)
		if err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
		for _, p := range rolePermissions {
			newRole.Permissions = append(newRole.Permissions, roles.PermissionVerbose{Path: p[1], Method: p[2]})
		}
		domainRoles = append(domainRoles, newRole)
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: domainRoles})
}
