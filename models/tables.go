package models

import (
	"github.com/scylladb/gocqlx/v2/table"
)

type TablesModel struct {
	InternalName   string                `cql:"internal_name"`
	Name           string                `cql:"name"`
	Description    string                `cql:"description"`
	PartitionKeys  []string              `cql:"partition_keys"`
	ClusteringKeys []string              `cql:"clustering_keys"`
	AllColumns     []string              `cql:"all_columns"`
	Mappings       map[string]string     `cql:"mappings"`
	Indexes        map[string]IndexModel `cql:"indexes"`
}

var TablesMetadata = table.Metadata{
	Name:    "Tables",
	Columns: []string{"internal_name", "name", "description", "partition_keys", "clustering_keys", "all_columns", "mappings"},
	PartKey: []string{"internal_name"},
	SortKey: []string{"name", "description"},
}

type IndexModel struct {
	Local     bool     `cql:"local"`
	IndexName string   `cql:"index_name"`
	TableName string   `cql:"table_name"`
	Columns   []string `cql:"columns"`
}
