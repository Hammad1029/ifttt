package controllers

import (
	"encoding/hex"
	"fmt"
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/common"
	"ifttt/manager/domain/user"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authController struct {
	serverCore *core.ServerCore
}

func newAuthController(serverCore *core.ServerCore) *authController {
	return &authController{
		serverCore: serverCore,
	}
}

func (a *authController) Login(c *gin.Context) {
	err, reqBodyAny := middlewares.Validator(c, user.LoginRequest{})
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	reqBody := reqBodyAny.(*user.LoginRequest)

	user, err := a.serverCore.ConfigStore.UserRepo.GetUser(reqBody.Email, user.DecodeUser)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	hashStr := hex.EncodeToString(hashBytes)
	if user == nil || user.Password != hashStr {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["WrongCredentials"]})
		return
	}

	tokenPair, err := a.serverCore.TokenService.NewTokenPair(reqBody.Email)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: tokenPair})
}

func (a *authController) Logout(c *gin.Context) {
	user := user.GetUserFromContext(c)
	if user == nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["UserNotFound"]})
		return
	}
	if err := a.serverCore.CacheStore.TokenRepo.DeleteTokenPair(user.Email); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	common.ResponseHandler(c, common.ResponseConfig{})
}

func (a *authController) RefreshToken(c *gin.Context) {
	tokenDetails, err := a.serverCore.TokenService.VerifyToken(c.GetHeader("refresh_token"))
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if tokenDetails == nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["Unauthorized"]})
		return
	}

	if float64(time.Now().Unix()) > float64(tokenDetails.Expiry) {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["Unauthorized"]})
		return
	}

	user, err := a.serverCore.ConfigStore.UserRepo.GetUser(
		fmt.Sprint(tokenDetails.Claims["email"]), user.DecodeUser,
	)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if user == nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["UserNotFound"]})
		return
	}

	if err := a.serverCore.CacheStore.TokenRepo.DeleteTokenPair(user.Email); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	tokenPair, err := a.serverCore.TokenService.NewTokenPair(user.Email)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: tokenPair})
}
