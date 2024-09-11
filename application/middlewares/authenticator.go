package middlewares

import (
	"fmt"
	"ifttt/manager/common"
	"ifttt/manager/domain/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (a *AllMiddlewares) Authenticator(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenDetails, err := a.serverCore.TokenService.VerifyToken(authHeader)
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

	userEmail := fmt.Sprint(tokenDetails.Claims["email"])

	ctx := c.Request.Context()
	cacheExists, err := a.serverCore.CacheStore.TokenRepo.GetTokenPair(userEmail, ctx)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if cacheExists == nil || !cacheExists.AccessToken.IsSameToken(authHeader) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := a.serverCore.ConfigStore.UserRepo.GetUser(userEmail, user.DecodeUser)
	if err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if user == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
	c.Set("subject", user.Email)
	c.Next()
}
