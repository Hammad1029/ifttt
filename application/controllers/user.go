package controllers

import (
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/common"
	"ifttt/manager/domain/user"

	"github.com/gin-gonic/gin"
)

type userController struct {
	serverCore *core.ServerCore
}

func newUserController(serverCore *core.ServerCore) *userController {
	return &userController{
		serverCore: serverCore,
	}
}

func (u *userController) CreateUser(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, user.CreateUserRequest{})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*user.CreateUserRequest)

	if user, err := u.serverCore.ConfigStore.UserRepo.GetUser(reqBody.Email, user.DecodeUser); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if user != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["UserAlreadyExists"]})
		return
	}

	newUser := user.User{Email: reqBody.Email, Password: reqBody.Password}
	if err := u.serverCore.ConfigStore.UserRepo.CreateUser(newUser); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{})
}

func (u *userController) GetAllUsers(c *gin.Context) {
	users, err := u.serverCore.ConfigStore.UserRepo.GetAllUsers()
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: users})
}
