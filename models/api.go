package models

import (
	"fmt"

	"github.com/scylladb/gocqlx/v2/table"
)

type ApiModel struct {
	ApiGroup       string              `cql:"api_group"`
	ApiName        string              `cql:"api_name"`
	ApiDescription string              `cql:"api_description"`
	ApiPath        string              `cql:"api_path"`
	ApiRequest     string              `cql:"api_request"`
	StartRules     []string            `cql:"start_rules"`
	Rules          map[string]*RuleUDT `cql:"rules"`
	Queries        map[string]QueryUDT `cql:"queries"`
}

var ApisMetadata = table.Metadata{
	Name:    "Apis",
	Columns: []string{"api_group", "api_name", "api_description", "api_path", "start_rules", "rules", "queries"},
	PartKey: []string{"api_group"},
	SortKey: []string{"api_name", "api_description"},
}

func (api *ApiModel) TransformApiForSave() error {
	api.Queries = make(map[string]QueryUDT)
	for _, rule := range api.Rules {
		if err := rule.TransformForSave(&api.Queries); err != nil {
			return fmt.Errorf("method TransformApiForSave: %s", err)
		}
	}
	return nil
}
