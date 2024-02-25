package models

type IndexModel struct {
	Local     bool
	Name      string
	TableName string
	Columns   []string
}
