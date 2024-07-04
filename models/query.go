package models

type QueryUDT struct {
	QueryString string          `cql:"query_str"`
	Resolvables []ResolvableUDT `cql:"resolvables"`
	Type        string          `cql:"type"`
}
