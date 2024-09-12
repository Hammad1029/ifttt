package controllers

import (
	"fmt"
	"ifttt/manager/application/core"
	"ifttt/manager/common"
	"ifttt/manager/domain/auth"
	"ifttt/manager/domain/user"
	"net/http"
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
	var reqBody auth.LoginRequest
	if ok := validateAndBind(c, &reqBody); !ok {
		return
	}

	user, err := a.serverCore.ConfigStore.UserRepo.GetUser(reqBody.Email, user.DecodeUser)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if user == nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["WrongCredentials"]})
		return
	}

	if pwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); pwErr != nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["WrongCredentials"]})
		return
	}

	tokenPair, err := a.serverCore.TokenService.NewTokenPair(reqBody.Email)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := a.serverCore.CacheStore.AuthRepo.DeleteTokenPair(reqBody.Email, ctx); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	if err := a.serverCore.CacheStore.AuthRepo.StoreTokenPair(reqBody.Email, tokenPair, ctx); err != nil {
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

	ctx := c.Request.Context()
	if err := a.serverCore.CacheStore.AuthRepo.DeleteTokenPair(user.Email, ctx); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	common.ResponseHandler(c, common.ResponseConfig{})
}

func (a *authController) RefreshToken(c *gin.Context) {
	refreshHeader := c.GetHeader("refresh_token")
	tokenDetails, err := a.serverCore.TokenService.VerifyToken(refreshHeader)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if tokenDetails == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > float64(tokenDetails.Expiry) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := c.Request.Context()
	userEmail := fmt.Sprint(tokenDetails.Claims["email"])

	cacheExists, err := a.serverCore.CacheStore.AuthRepo.GetTokenPair(userEmail, ctx)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if cacheExists == nil || !cacheExists.RefreshToken.IsSameToken(refreshHeader) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := a.serverCore.ConfigStore.UserRepo.GetUser(userEmail, user.DecodeUser)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if user == nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["UserNotFound"]})
		return
	}

	if err := a.serverCore.CacheStore.AuthRepo.DeleteTokenPair(user.Email, ctx); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	tokenPair, err := a.serverCore.TokenService.NewTokenPair(user.Email)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	}
	if err := a.serverCore.CacheStore.AuthRepo.StoreTokenPair(user.Email, tokenPair, ctx); err != nil {
		common.HandleErrorResponse(c, err)
		return
	}

	common.ResponseHandler(c, common.ResponseConfig{Data: tokenPair})
}
