package server

import (
	"fmt"
	"ifttt/manager/application/config"
	"ifttt/manager/application/controllers"
	"ifttt/manager/application/middlewares"
	infrastructure "ifttt/manager/infrastructure/init"

	"github.com/gin-gonic/gin"
)

func Init() error {
	dbStore, err := infrastructure.NewDbStore()
	if err != nil {
		return fmt.Errorf("method Init: error in creating db store: %s", err)
	}

	controllers := controllers.NewAllController(dbStore)

	port := config.GetConfigProp("app.port")
	router := gin.New()
	router.Use(middlewares.CORSMiddleware())
	plugRoutes(router, controllers)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		return fmt.Errorf("method Init: error in running gin router: %s", err)
	}
	return nil
}
