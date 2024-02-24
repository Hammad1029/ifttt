package routes

import (
	"generic/controllers"

	"github.com/gin-gonic/gin"
)

func apiRoutes(group *gin.RouterGroup) {
	group.POST("/addApi", controllers.AddApi)
	group.POST("/addMappingToApi", controllers.AddMappingToApi)
	group.POST("/Call/:clientId/:pathName", controllers.CallApi)
}
