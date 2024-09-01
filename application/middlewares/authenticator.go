package middlewares

import (
	"fmt"
	"ifttt/manager/domain/user"
	"ifttt/manager/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (a *allMiddlewares) Authenticator(c *gin.Context) {
	tokenDetails, err := a.serverCore.TokenService.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	} else if tokenDetails == nil {
		utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["Unauthorized"]})
		return
	}

	if float64(time.Now().Unix()) > float64(tokenDetails.Expiry) {
		utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["Unauthorized"]})
		return
	}

	userEmail := fmt.Sprint(tokenDetails.Claims["email"])

	cacheExists, err := a.serverCore.CacheStore.TokenRepo.GetTokenPair(userEmail)
	if cacheExists == nil {
		utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["Unauthorized"]})
		return
	} else if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	}

	user, err := a.serverCore.ConfigStore.UserRepo.GetUser(userEmail, user.DecodeUser)
	if err != nil {
		utils.HandleErrorResponse(c, err)
		return
	} else if user == nil {
		utils.ResponseHandler(c, utils.ResponseConfig{Response: utils.Responses["UserNotFound"]})
		return
	}

	c.Set("user", user)
	c.Next()
}
