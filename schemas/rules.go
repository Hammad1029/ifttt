package schemas

import "generic/models"

type AddRuleToApi struct {
	ClientId string                `json:"clientId" binding:"required"`
	Rule     models.RuleModelMongo `json:"rule" binding:"required"`
	ApiName  string                `json:"apiName" binding:"required"`
}
