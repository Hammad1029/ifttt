package models

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type RuleUDT struct {
	Id          string       `cql:"id" json:"id"`
	Description string       `cql:"description" json:"description"`
	Conditions  Condition    `cql:"conditions" json:"conditions"`
	Then        []Resolvable `cql:"then" json:"then"`
	Else        []Resolvable `cql:"else" json:"else"`
}

type Condition struct {
	ConditionType string      `cql:"condition_type" json:"conditionType" mapstructure:"conditionType"`
	Conditions    []Condition `cql:"conditions" json:"conditions" mapstructure:"conditions"`
	Group         bool        `cql:"group" json:"group" mapstructure:"group"`
	Operator1     Resolvable  `cql:"op1" json:"op1" mapstructure:"op1"`
	Operand       string      `cql:"opnd" json:"opnd" mapstructure:"opnd"`
	Operator2     Resolvable  `cql:"op2" json:"op2" mapstructure:"op2"`
}

func (r *RuleUDT) TransformForSave(queries *map[string]QueryUDT) error {
	cGroup := Condition{}
	if err := mapstructure.Decode(r.Conditions, &cGroup); err != nil {
		return fmt.Errorf("method TransformForSave: %s", err)
	}
	cGroup.transformGroup(queries)

	for _, ac := range r.Then {
		if err := ac.generateQueries(queries); err != nil {
			return err
		}
	}

	for _, ac := range r.Else {
		if err := ac.generateQueries(queries); err != nil {
			return err
		}
	}

	return nil
}

func (c *Condition) transformGroup(queries *map[string]QueryUDT) error {
	for _, cond := range c.Conditions {
		if cond.Group {
			cond.transformGroup(queries)
		} else {
			if err := cond.Operator1.generateQueries(queries); err != nil {
				return fmt.Errorf("method transformGroup: %s", err)
			}
			if err := cond.Operator2.generateQueries(queries); err != nil {
				return fmt.Errorf("method transformGroup: %s", err)
			}
		}
	}
	return nil
}
