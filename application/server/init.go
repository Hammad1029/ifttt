package server

import (
	"fmt"
	"ifttt/manager/application/config"
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"

	"github.com/gin-gonic/gin"
)

var serverCore *core.ServerCore

func Init() error {
	if core, err := core.NewServerCore(); err != nil {
		return fmt.Errorf("could not create server core: %s", err)
	} else {
		serverCore = core
	}

	port := config.GetConfigProp("app.port")
	router := gin.New()
	router.Use(middlewares.CORSMiddleware())
	plugRoutes(router, serverCore)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		return fmt.Errorf("method Init: error in running gin router: %s", err)
	}
	return nil
}
