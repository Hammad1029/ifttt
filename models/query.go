package models

type QueryUDT struct {
	QueryString string       `cql:"query_str"`
	Resolvables []Resolvable `cql:"resolvables"`
	Type        string       `cql:"type"`
}
