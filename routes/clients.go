package routes

import (
	"generic/controllers"

	"github.com/gin-gonic/gin"
)

func clientRoutes(group *gin.RouterGroup) {
	group.POST("/addClient", controllers.AddClient)
	group.POST("/addApi", controllers.AddApi)
}
