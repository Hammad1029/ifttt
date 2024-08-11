package schemas

import (
	"generic/domain/api"
)

type CreateApiRequest struct {
	ApiGroup       string               `json:"apiGroup"`
	ApiName        string               `json:"apiName"`
	ApiDescription string               `json:"apiDescription"`
	ApiPath        string               `json:"apiPath"`
	ApiRequest     map[string]any       `json:"apiRequest"`
	StartRules     []string             `json:"startRules"`
	Rules          map[string]*api.Rule `json:"rules"`
}

type GetApisRequest struct {
	ApiGroup       string `cql:"apiGroup" json:"apiGroup"`
	ApiName        string `cql:"apiName" json:"apiName"`
	ApiDescription string `cql:"apiDescription" json:"apiDescription"`
}
