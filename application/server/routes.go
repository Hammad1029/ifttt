package server

import (
	"ifttt/manager/application/controllers"

	"github.com/gin-gonic/gin"
)

func plugRoutes(router *gin.Engine, controllers *controllers.AllController) {
	apisGroup := router.Group("/apis")
	apisGroup.POST("/createApi", controllers.ApiController.CreateApi)
	apisGroup.GET("/getApis", controllers.ApiController.GetApis)
}
