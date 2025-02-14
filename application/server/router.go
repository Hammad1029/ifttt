package server

import (
	"ifttt/manager/application/controllers"
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/common"
	"ifttt/manager/domain/roles"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initRouter(router *gin.Engine, serverCore *core.ServerCore) {
	controllers := controllers.NewAllController(serverCore)
	serverCore.Routes = getRouteDefinitions(controllers)
	serverCore.Permissions = &[]string{}
	middlewares := middlewares.NewAllMiddlewares(serverCore)

	router.Use(middlewares.CORS())
	routerGroup := router.Group("/")
	createRoutes(routerGroup, middlewares,
		serverCore.Routes, serverCore.Permissions)
}

func createRoutes(
	routerGroup *gin.RouterGroup, middlewares *middlewares.AllMiddlewares,
	routeDefintions *[]common.RouteDefinition, permissions *[]string) {
	for _, r := range *routeDefintions {
		if r.Method == "GROUP" {
			newGroup := routerGroup.Group(r.Path)
			createRoutes(newGroup, middlewares, &r.Children, permissions)
			continue
		}

		handlers := []gin.HandlerFunc{}
		if r.Authenticated {
			handlers = append(handlers, middlewares.Authenticator)
		}
		if r.Authorized {
			newPermission := roles.PermissionVerbose{Path: routerGroup.BasePath() + r.Path, Method: r.Method}
			*permissions = append(*permissions, newPermission.CreatePermission())
			handlers = append(handlers, middlewares.Authorizer)
		}
		handlers = append(handlers, r.HandlerFunc)

		switch r.Method {
		case http.MethodGet:
			routerGroup.GET(r.Path, handlers...)
		case http.MethodPost:
			routerGroup.POST(r.Path, handlers...)
		}
	}
}
