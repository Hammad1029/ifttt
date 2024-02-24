package models

import (
	"github.com/scylladb/gocqlx/v2/table"
)

type TablesModel struct {
	InternalName   string
	Name           string
	Description    string
	PartitionKeys  []string
	ClusteringKeys []string
	AllColumns     []string
	Mappings       map[string]string
}

var TablesMetadata = table.Metadata{
	Name:    "Tables",
	Columns: []string{"internal_name", "name", "description", "partition_keys", "clustering_keys", "all_columns", "mappings"},
	PartKey: []string{"internal_name"},
	SortKey: []string{"name", "description"},
}
