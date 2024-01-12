package routes

import (
	"generic/controllers"
	"generic/schemas"
	"generic/utils"

	"github.com/gin-gonic/gin"
)

func rulesRoutes(group *gin.RouterGroup) {
	group.POST("/AddRuleToApi", utils.Validator(controllers.AddRuleToApi, utils.Roles["unauthorized"], schemas.AddRuleToApi{}))
	group.POST("/Call/:id", utils.Validator(controllers.RulesCall, utils.Roles["unauthorized"], nil))
}
