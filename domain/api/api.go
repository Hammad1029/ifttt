package api

import (
	"encoding/json"
	"fmt"

	"github.com/scylladb/gocqlx/v3/table"
)

type ApiSerialized struct {
	ApiGroup       string   `json:"apiGroup" mapstructure:"apiGroup"`
	ApiName        string   `json:"apiName" mapstructure:"apiName"`
	ApiDescription string   `json:"apiDescription" mapstructure:"apiDescription"`
	ApiPath        string   `json:"apiPath" mapstructure:"apiPath"`
	ApiRequest     string   `json:"apiRequest" mapstructure:"apiRequest"`
	StartRules     []string `json:"rules" mapstructure:"rules"`
	Rules          string   `json:"startRules" mapstructure:"startRules"`
}

type Api struct {
	ApiGroup       string           `json:"apiGroup" mapstructure:"apiGroup"`
	ApiName        string           `json:"apiName" mapstructure:"apiName"`
	ApiDescription string           `json:"apiDescription" mapstructure:"apiDescription"`
	ApiPath        string           `json:"apiPath" mapstructure:"apiPath"`
	ApiRequest     map[string]any   `json:"apiRequest" mapstructure:"apiRequest"`
	StartRules     []string         `json:"startRules" mapstructure:"startRules"`
	Rules          map[string]*Rule `json:"rules" mapstructure:"rules"`
}

var ApisMetadata = table.Metadata{
	Name:    "Apis",
	Columns: []string{"api_group", "api_name", "api_description", "api_path", "api_request", "start_rules", "rules", "queries"},
	PartKey: []string{"api_group"},
	SortKey: []string{"api_name", "api_description"},
}

func (a *Api) Serialize() (*ApiSerialized, error) {
	apiModelSerialized := ApiSerialized{
		ApiGroup:       a.ApiGroup,
		ApiName:        a.ApiName,
		ApiDescription: a.ApiDescription,
		ApiPath:        a.ApiPath,
		StartRules:     a.StartRules,
	}

	requestSerialized, err := json.Marshal(a.ApiRequest)
	if err != nil {
		return nil, fmt.Errorf("method *Api.TransformApiForSave: could not serialize request: %s", err)
	}
	apiModelSerialized.ApiRequest = string(requestSerialized)

	rulesSerialized, err := json.Marshal(a.Rules)
	if err != nil {
		return nil, fmt.Errorf("method *Api.TransformApiForSave: could not serialize rules: %s", err)
	}
	apiModelSerialized.Rules = string(rulesSerialized)

	return &apiModelSerialized, nil
}
