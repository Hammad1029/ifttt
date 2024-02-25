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
		// rulesRoutes(v1.Group("/rules"))
		// clientRoutes(v1.Group("/clients"))
		// apiRoutes(v1.Group("/apis"))
	}
	router.Run()
}
