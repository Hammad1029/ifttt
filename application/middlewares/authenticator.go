package middlewares

import (
	"fmt"
	"ifttt/manager/common"
	"ifttt/manager/domain/user"
	"time"

	"github.com/gin-gonic/gin"
)

func (a *allMiddlewares) Authenticator(c *gin.Context) {
	tokenDetails, err := a.serverCore.TokenService.VerifyToken(c.GetHeader("Authorization"))
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

	userEmail := fmt.Sprint(tokenDetails.Claims["email"])

	cacheExists, err := a.serverCore.CacheStore.TokenRepo.GetTokenPair(userEmail)
	if cacheExists == nil {
		common.ResponseHandler(c, common.ResponseConfig{Response: common.Responses["Unauthorized"]})
		return
	} else if err != nil {
		common.HandleErrorResponse(c, err)
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

	c.Set("user", user)
	c.Next()
}
