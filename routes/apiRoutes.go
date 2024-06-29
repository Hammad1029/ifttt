package routes

import (
	"generic/controllers"

	"github.com/gin-gonic/gin"
)

func apiRoutes(group *gin.RouterGroup) {
	group.POST("/addApi", controllers.Apis.AddApi)
	group.POST("/getApis", controllers.Apis.GetApis)
}
