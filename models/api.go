package models

import "github.com/scylladb/gocqlx/v2/table"

type ApiModel struct {
	ApiGroup       string              `cql:"api_group"`
	ApiName        string              `cql:"api_name"`
	ApiDescription string              `cql:"api_description"`
	ApiPath        string              `cql:"api_path"`
	StartRules     []int               `cql:"start_rules"`
	Rules          []RuleUDT           `cql:"rules"`
	Queries        map[string]QueryUDT `cql:"queries"`
}

var ApisMetadata = table.Metadata{
	Name:    "Apis",
	Columns: []string{"api_group", "api_name", "api_description", "api_path", "start_rules", "rules", "queries"},
	PartKey: []string{"api_group"},
	SortKey: []string{"api_name", "api_description"},
}

func (api *ApiModel) TransformApiForSave() {
	api.Queries = make(map[string]QueryUDT)
	for _, rule := range api.Rules {
		rule.TransformForSave(&api.Queries)
	}
}
