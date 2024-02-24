package server

import (
	"generic/routes"

	"github.com/gin-gonic/gin"
)

func Init() {
	// port := config.GetConfigProp("app.port")
	router := gin.New()
	routes.Init(router)
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	router.Run(":" + "5600")
}
