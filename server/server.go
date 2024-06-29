package server

import (
	"generic/config"
	"generic/middlewares"
	"generic/routes"

	"github.com/gin-gonic/gin"
)

func Init() {
	port := config.GetConfigProp("app.port")
	router := gin.New()
	router.Use(middlewares.CORSMiddleware())
	routes.Init(router)
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	router.Run(":" + port)
}
