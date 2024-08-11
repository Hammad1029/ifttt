package server

import (
	"fmt"
	"generic/application/config"
	"generic/application/controllers"
	"generic/application/middlewares"
	infrastructure "generic/infrastructure/init"

	"github.com/gin-gonic/gin"
)

func Init() {
	dbStore, err := infrastructure.NewDbStore()
	if err != nil {
		panic(fmt.Errorf("error in creating db store: %s", err))
	}

	controllers := controllers.NewAllController(dbStore)

	port := config.GetConfigProp("app.port")
	router := gin.New()
	router.Use(middlewares.CORSMiddleware())
	plugRoutes(router, controllers)

	router.Run(fmt.Sprintf(":%s", port))
}
