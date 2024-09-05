package server

import (
	"ifttt/manager/application/controllers"
	"ifttt/manager/application/core"
	"ifttt/manager/application/middlewares"
	"ifttt/manager/common"

	"github.com/gin-contrib/authz"
	"github.com/gin-gonic/gin"
)

func plugRoutes(router *gin.Engine, serverCore *core.ServerCore) {
	controllers := controllers.NewAllController(serverCore)
	middlewares := middlewares.NewAllMiddlewares(serverCore)

	router.Use(middlewares.CORS())
	router.Use(authz.NewAuthorizer(serverCore.ConfigStore.CasbinEnforcer))

	authRouter := router.Group(common.GetRoutePath(common.AuthGroup))
	authRouter.POST(common.GetRoutePath(common.AuthLogin), controllers.AuthController.Login)
	authRouter.POST(common.GetRoutePath(common.AuthRefresh), controllers.AuthController.RefreshToken)
	router.Use(middlewares.Authenticator)
	authRouter.GET(common.GetRoutePath(common.AuthLogout), controllers.AuthController.Logout)

	apisRouter := router.Group(common.GetRoutePath(common.ApisGroup))
	apisRouter.GET(common.GetRoutePath(common.ApisGet), controllers.ApiController.GetApis)
	apisRouter.POST(common.GetRoutePath(common.ApisCreate), controllers.ApiController.CreateApi)

	schemasRouter := router.Group(common.GetRoutePath(common.SchemasGroup))
	schemasRouter.GET(common.GetRoutePath(common.SchemasGet), controllers.TablesController.GetSchema)
	schemasRouter.POST(common.GetRoutePath(common.SchemasCreateTable), controllers.TablesController.CreateTable)
	schemasRouter.POST(common.GetRoutePath(common.SchemasUpdateTable), controllers.TablesController.UpdateTable)

	userRouter := router.Group(common.GetRoutePath(common.UserGroup))
	userRouter.POST(common.GetRoutePath(common.UserCreate), controllers.UserController.CreateUser)

	rolesRouter := router.Group(common.GetRoutePath(common.RolesGroup))
	rolesRouter.GET(common.GetRoutePath(common.RolesGetAllPermissions), controllers.RoleController.GetAllPermissions)
	rolesRouter.GET(common.GetRoutePath(common.RolesAddUpdateRole), controllers.RoleController.AddUpdateRole)
	rolesRouter.GET(common.GetRoutePath(common.RolesUpdateUserRoles), controllers.RoleController.UpdateUserRoles)
}
