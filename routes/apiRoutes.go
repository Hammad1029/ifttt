package routes

import (
	"generic/controllers"
	"generic/schemas"
	"generic/utils"

	"github.com/gin-gonic/gin"
)

func apiRoutes(group *gin.RouterGroup) {
	group.POST("/addApi", utils.Validator(controllers.AddApi, utils.Roles["unauthorized"], schemas.AddApi{}))
	group.POST("/addMappingToApi", utils.Validator(controllers.AddMappingToApi, utils.Roles["unauthorized"], schemas.AddMappingToApi{}))
	group.POST("/Call/:clientId/:pathName", utils.Validator(controllers.CallApi, utils.Roles["unauthorized"], nil))
}
