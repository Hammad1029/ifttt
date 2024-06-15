package models

import (
	"github.com/scylladb/gocqlx/v2/table"
)

type TableColumn struct {
	Name          string `cql:"name" json:"name"`
	DataType      string `cql:"data_type" json:"dataType"`
	PartitionKey  bool   `cql:"partition_key" json:"partitionKey"`
	ClusteringKey bool   `cql:"clustering_key" json:"clusteringKey"`
}

type TablesModel struct {
	InternalName string                `cql:"internal_name" json:"internalName"`
	Name         string                `cql:"name" json:"name"`
	Description  string                `cql:"description" json:"description"`
	Columns      []TableColumn         `cql:"columns" json:"columns"`
	Indexes      map[string]IndexModel `cql:"indexes" json:"indexes"`
}

var TablesMetadata = table.Metadata{
	Name:    "Tables",
	Columns: []string{"internal_name", "name", "description", "columns", "indexes"},
	PartKey: []string{"internal_name"},
	SortKey: []string{"name", "description"},
}

type IndexModel struct {
	Local     bool     `cql:"local"`
	IndexName string   `cql:"index_name"`
	TableName string   `cql:"table_name"`
	Columns   []string `cql:"columns"`
}
