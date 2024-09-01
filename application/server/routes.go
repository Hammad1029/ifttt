package server

import (
	"ifttt/manager/application/controllers"
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"

	"github.com/gin-gonic/gin"
)

func plugRoutes(router *gin.Engine, serverCore *core.ServerCore) {
	controllers := controllers.NewAllController(serverCore)
	middlewares := middlewares.NewAllMiddlewares(serverCore)

	router.Use(middlewares.CORS())

	authGroup := router.Group("/auth")
	authGroup.POST("/login", controllers.AuthController.Login)
	authGroup.POST("/refresh", controllers.AuthController.RefreshToken)
	router.Use(middlewares.Authenticator)
	authGroup.GET("/logout", controllers.AuthController.Logout)

	apisGroup := router.Group("/apis")
	apisGroup.GET("/", controllers.ApiController.GetApis)
	apisGroup.POST("/create", controllers.ApiController.CreateApi)

	schemasGroup := router.Group("/schemas")
	schemasGroup.GET("/", controllers.TablesController.GetSchema)
	schemasGroup.POST("/create", controllers.TablesController.CreateTable)
	schemasGroup.POST("/update", controllers.TablesController.UpdateTable)
}
