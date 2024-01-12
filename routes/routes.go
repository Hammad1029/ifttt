package routes

import "github.com/gin-gonic/gin"

func Init(router *gin.Engine) {
	v1 := router.Group("/api")
	{
		rulesRoutes(v1.Group("/rules"))
		clientRoutes(v1.Group("/clients"))
		apiRoutes(v1.Group("/apis"))
	}
	router.Run()
}
