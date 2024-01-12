package server

import (
	"fmt"
	"generic/config"
	"generic/routes"

	"github.com/gin-gonic/gin"
)

func Init() {
	port := config.GetConfigProp("app.port")
	router := gin.New()
	routes.Init(router)
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	router.Run(":" + port)
	fmt.Println("Server started at port: ", port)
}
