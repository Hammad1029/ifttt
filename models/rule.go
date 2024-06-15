package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"generic/utils"

	jsontocql "github.com/Hammad1029/json-to-cql"
	"github.com/mitchellh/mapstructure"
)

type RuleUDT struct {
	Conditions interface{}     `cql:"conditions" json:"conditions"`
	Then       []ResolvableUDT `cql:"then" json:"then"`
	Else       []ResolvableUDT `cql:"else" json:"else"`
}

type Condition struct {
	ConditionType string        `cql:"condition_type" json:"conditionType" mapstructure:"conditionType"`
	Conditions    []Condition   `cql:"conditions" json:"conditions" mapstructure:"conditions"`
	Group         bool          `cql:"group" json:"group" mapstructure:"group"`
	Operator1     ResolvableUDT `cql:"op1" json:"op1" mapstructure:"op1"`
	Operand       string        `cql:"opnd" json:"opnd" mapstructure:"opnd"`
	Operator2     ResolvableUDT `cql:"op2" json:"op2" mapstructure:"op2"`
}

type ResolvableUDT struct {
	Type string                 `cql:"type" json:"type" mapstructure:"type"`
	Data map[string]interface{} `cql:"data" json:"data" mapstructure:"data"`
}

func (r *RuleUDT) TransformForSave(queries *map[string]QueryUDT) error {
	cGroup := Condition{}
	if err := mapstructure.Decode(r.Conditions, &cGroup); err != nil {
		return fmt.Errorf("method TransformForSave: %s", err)
	}
	cGroup.transformGroup(queries)

	marshalled, err := json.Marshal(cGroup)
	if err != nil {
		return fmt.Errorf("method TransformForSave: %s", err)
	}
	r.Conditions = string(marshalled)

	for _, ac := range r.Then {
		if err := ac.generateQueries(queries); err != nil {
			return err
		} else {
			if stringified, err := utils.StringifyMapInt(ac.Data); err != nil {
				return err
			} else {
				ac.Data = stringified
			}
		}
	}

	for _, ac := range r.Else {
		if err := ac.generateQueries(queries); err != nil {
			return err
		} else {
			if stringified, err := utils.StringifyMapInt(ac.Data); err != nil {
				return err
			} else {
				ac.Data = stringified
			}
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

func (ruleResolvable *ResolvableUDT) generateQueries(queries *map[string]QueryUDT) error {
	if query, ok := ruleResolvable.Data["query"]; ruleResolvable.Type == "db" && ok {
		if queryMap, ok := query.(map[string]interface{}); ok {
			var queryDoc jsontocql.QueryDoc
			mapstructure.Decode(queryMap, &queryDoc)
			if paraQuery, err := queryDoc.CreateParameterizedQuery(); err != nil {
				return err
			} else {
				ruleResolvable.Data["query"] = paraQuery.QueryHash

				var queryResolvables []ResolvableUDT
				mapstructure.Decode(paraQuery.Resolvables, &queryResolvables)

				(*queries)[paraQuery.QueryHash] = QueryUDT{
					QueryString: paraQuery.QueryString,
					Resolvables: queryResolvables,
					Type:        paraQuery.Type,
				}
				return nil
			}
		} else {
			return errors.New("failure in typecasting map")
		}
	} else {
		return nil
	}
}
