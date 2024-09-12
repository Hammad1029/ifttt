package middlewares

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *AllMiddlewares) Authorizer(c *gin.Context) {
	sub := user.GetUserFromContext(c).Email
	obj := c.Request.URL.Path
	act := c.Request.Method
	if ok, err := a.serverCore.ConfigStore.CasbinEnforcer.Enforce(sub, obj, act); err != nil {
		common.HandleErrorResponse(c, err)
		return
	} else if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
