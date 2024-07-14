package models

import (
	"encoding/json"
	"fmt"

	jsontocql "github.com/Hammad1029/json-to-cql"
	"github.com/scylladb/gocqlx/v2/table"
)

type ApiModelSerialized struct {
	ApiGroup       string   `cql:"api_group"`
	ApiName        string   `cql:"api_name"`
	ApiDescription string   `cql:"api_description"`
	ApiPath        string   `cql:"api_path"`
	ApiRequest     string   `cql:"api_request"`
	StartRules     []string `cql:"start_rules"`
	Rules          string   `cql:"rules"`
	Queries        string   `cql:"queries"`
}

type ApiModel struct {
	ApiGroup       string                                  `json:"api_group" mapstructure:"api_group"`
	ApiName        string                                  `json:"api_name" mapstructure:"api_name"`
	ApiDescription string                                  `json:"api_description" mapstructure:"api_description"`
	ApiPath        string                                  `json:"api_path" mapstructure:"api_path"`
	ApiRequest     map[string]interface{}                  `json:"api_request" mapstructure:"api_request"`
	StartRules     []string                                `json:"start_rules" mapstructure:"start_rules"`
	Rules          map[string]*RuleUDT                     `json:"rules" mapstructure:"rules"`
	Queries        map[string]jsontocql.ParameterizedQuery `json:"queries" mapstructure:"queries"`
}

var ApisMetadata = table.Metadata{
	Name:    "Apis",
	Columns: []string{"api_group", "api_name", "api_description", "api_path", "api_request", "start_rules", "rules", "queries"},
	PartKey: []string{"api_group"},
	SortKey: []string{"api_name", "api_description"},
}

func (api *ApiModel) TransformApiForSave() (ApiModelSerialized, error) {
	api.Queries = make(map[string]jsontocql.ParameterizedQuery)
	for _, rule := range api.Rules {
		if err := rule.TransformForSave(&api.Queries); err != nil {
			return ApiModelSerialized{}, fmt.Errorf("method TransformApiForSave: %s", err)
		}
	}

	requestSerialized, err := json.Marshal(api.ApiRequest)
	if err != nil {
		return ApiModelSerialized{}, fmt.Errorf("method TransformApiForSave: %s", err)
	}

	rulesSerialized, err := json.Marshal(api.Rules)
	if err != nil {
		return ApiModelSerialized{}, fmt.Errorf("method TransformApiForSave: %s", err)
	}

	queriesSerialized, err := json.Marshal(api.Queries)
	if err != nil {
		return ApiModelSerialized{}, fmt.Errorf("method TransformApiForSave: %s", err)
	}

	return ApiModelSerialized{
		ApiGroup:       api.ApiGroup,
		ApiName:        api.ApiName,
		ApiDescription: api.ApiDescription,
		ApiPath:        api.ApiPath,
		ApiRequest:     string(requestSerialized),
		StartRules:     api.StartRules,
		Rules:          string(rulesSerialized),
		Queries:        string(queriesSerialized),
	}, nil
}
