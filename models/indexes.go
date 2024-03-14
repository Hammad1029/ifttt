package models

type IndexModel struct {
	Local     bool     `cql:"local"`
	IndexName string   `cql:"index_name"`
	TableName string   `cql:"table_name"`
	Columns   []string `cql:"columns"`
}
