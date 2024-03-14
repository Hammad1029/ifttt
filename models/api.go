package models

import "github.com/scylladb/gocqlx/v2/table"

type ApiModel struct {
	ApiGroup       string    `cql:"api_group"`
	ApiName        string    `cql:"api_name"`
	ApiDescription string    `cql:"api_description"`
	ApiPath        string    `cql:"api_path"`
	Rules          []RuleUDT `cql:"rules"`
}

var ApisMetadata = table.Metadata{
	Name:    "Apis",
	Columns: []string{"api_group", "api_name", "api_description", "api_path", "rules"},
	PartKey: []string{"api_group"},
	SortKey: []string{"api_name", "api_description"},
}
