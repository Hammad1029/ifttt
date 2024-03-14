package models

import (
	"github.com/gocql/gocql"
)

type RuleUDT struct {
	Id        gocql.UUID  `cql:"id"`
	Operator1 string      `cql:"op1"`
	Operand   string      `cql:"opnd"`
	Operator2 string      `cql:"op2"`
	Then      []ActionUDT `cql:"then"`
	Else      []ActionUDT `cql:"else"`
}

type ActionUDT struct {
	Type string `cql:"type"`
	Data string `cql:"data"`
}
