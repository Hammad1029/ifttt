package routes

import (
	"github.com/gin-gonic/gin"
)

func rulesRoutes(group *gin.RouterGroup) {
	group.POST("/AddRuleToApi")
	group.POST("/Call/:id")
}
