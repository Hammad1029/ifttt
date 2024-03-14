package routes

import (
	"generic/controllers"

	"github.com/gin-gonic/gin"
)

func indexesRoutes(group *gin.RouterGroup) {
	group.POST("/addIndex", controllers.Indexes.AddIndex)
	group.POST("/findIndex", controllers.Indexes.FindIndex)
	group.POST("/dropIndex", controllers.Indexes.DropIndex)
}
