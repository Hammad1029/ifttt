package server

import (
	"ifttt/manager/application/controllers"
	"ifttt/manager/application/core"

	"github.com/gin-gonic/gin"
)

func plugRoutes(router *gin.Engine, serverCore *core.ServerCore) {
	controllers := controllers.NewAllController(serverCore)

	apisGroup := router.Group("/apis")
	apisGroup.GET("/", controllers.ApiController.GetApis)
	apisGroup.POST("/create", controllers.ApiController.CreateApi)

	schemasGroup := router.Group("/schemas")
	schemasGroup.GET("/", controllers.TablesController.GetSchema)
	schemasGroup.POST("/create", controllers.TablesController.CreateTable)
	schemasGroup.POST("/update", controllers.TablesController.UpdateTable)
}
