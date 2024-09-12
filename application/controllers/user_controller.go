package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/roles"
	"ifttt/manager/domain/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userController struct {
	serverCore *core.ServerCore
}

func newUserController(serverCore *core.ServerCore) *userController {
	return &userController{
		serverCore: serverCore,
	}
}

func (uc *userController) CreateUser(c *gin.Context) {
	var reqBody user.CreateUserRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	if user, err := uc.serverCore.ConfigStore.UserRepo.GetUser(reqBody.Email, user.DecodeUser); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if user != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["UserAlreadyExists"]})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	newUser := user.User{Email: reqBody.Email, Password: string(hashedPassword)}
	if err := uc.serverCore.ConfigStore.UserRepo.CreateUser(newUser); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (uc *userController) GetAllUsers(c *gin.Context) {
	users, err := uc.serverCore.ConfigStore.UserRepo.GetAllUsers()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	for _, u := range users {
		userRoles, err := uc.serverCore.ConfigStore.CasbinEnforcer.GetRolesForUser(u.Email)
		if err != nil {
			common.HandleErrorResponse(c, err)
			return
		}
		for _, r := range userRoles {
			newRole := roles.Role{RoleName: r}
			rolePermissions, err := uc.serverCore.ConfigStore.CasbinEnforcer.GetPermissionsForUser(r)
			if err != nil {
				common.HandleErrorResponse(c, err)
				return
			}
			for _, p := range rolePermissions {
				newRole.Permissions = append(newRole.Permissions, roles.PermissionVerbose{Path: p[1], Method: p[2]})
			}
			u.Roles = append(u.Roles, newRole)
		}
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: users})
}
