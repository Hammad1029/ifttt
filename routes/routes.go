package routes

import (
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	// middlewares.ValidatorInit()
	v1 := router.Group("/api")
	{
		tablesRoutes(v1.Group("/tables"))
		indexesRoutes(v1.Group("/indexes"))
		apiRoutes(v1.Group("/apis"))
		// rulesRoutes(v1.Group("/rules"))
		// clientRoutes(v1.Group("/clients"))
	}
	router.Run()
}
