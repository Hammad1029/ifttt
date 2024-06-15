package routes

import (
	"generic/controllers"

	"github.com/gin-gonic/gin"
)

func tablesRoutes(group *gin.RouterGroup) {
	group.POST("/addTable", controllers.Tables.AddTable)
	group.POST("/getTables", controllers.Tables.GetTables)
	group.POST("/updateTable", controllers.Tables.UpdateTable)
}
