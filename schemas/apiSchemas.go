package schemas

import (
	"generic/models"
)

type AddApiRequest struct {
	ApiGroup       string                     `json:"apiGroup"`
	ApiName        string                     `json:"apiName"`
	ApiDescription string                     `json:"apiDescription"`
	ApiPath        string                     `json:"apiPath"`
	ApiRequest     string                     `json:"apiRequest"`
	StartRules     []string                   `json:"startRules"`
	Rules          map[string]*models.RuleUDT `json:"rules"`
}

type GetApisRequest struct {
	ApiGroup       string `cql:"apiGroup" json:"apiGroup"`
	ApiName        string `cql:"apiName" json:"apiName"`
	ApiDescription string `cql:"apiDescription" json:"apiDescription"`
}
