package models

type IndexModel struct {
	Local     bool     `cql:"local"`
	Name      string   `cql:"name"`
	TableName string   `cql:"table_name"`
	Columns   []string `cql:"columns"`
}
