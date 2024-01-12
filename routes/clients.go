package routes

import (
	"generic/controllers"
	"generic/schemas"
	"generic/utils"

	"github.com/gin-gonic/gin"
)

func clientRoutes(group *gin.RouterGroup) {
	group.POST("/addClient", utils.Validator(controllers.AddClient, utils.Roles["unauthorized"], schemas.AddClient{}))
	group.POST("/addApi", utils.Validator(controllers.AddApi, utils.Roles["unauthorized"], schemas.AddApi{}))
}
