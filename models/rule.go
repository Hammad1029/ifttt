package models

import (
	"errors"
	"generic/utils"

	jsontocql "github.com/Hammad1029/json-to-cql"
	"github.com/mitchellh/mapstructure"
)

type RuleUDT struct {
	Operator1 ResolvableUDT   `cql:"op1" json:"op1"`
	Operand   string          `cql:"opnd" json:"opnd"`
	Operator2 ResolvableUDT   `cql:"op2" json:"op2"`
	Then      []ResolvableUDT `cql:"then" json:"then"`
	Else      []ResolvableUDT `cql:"else" json:"else"`
}

type ResolvableUDT struct {
	Type string                 `cql:"type" json:"type"`
	Data map[string]interface{} `cql:"data" json:"data"`
}

type QueryUDT struct {
	QueryString string          `cql:"query_str"`
	Resolvables []ResolvableUDT `cql:"resolvables"`
	Type        string          `cql:"type"`
}

func (r *RuleUDT) TransformForSave(queries *map[string]QueryUDT) error {
	if err := r.generateQueries(&r.Operator1, queries); err != nil {
		return err
	} else {
		if stringified, err := utils.StringifyMapInt(r.Operator1.Data); err != nil {
			return err
		} else {
			r.Operator1.Data = stringified
		}
	}

	if err := r.generateQueries(&r.Operator2, queries); err != nil {
		return err
	} else {
		if stringified, err := utils.StringifyMapInt(r.Operator2.Data); err != nil {
			return err
		} else {
			r.Operator2.Data = stringified
		}
	}

	for _, ac := range r.Then {
		if err := r.generateQueries(&ac, queries); err != nil {
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
		if err := r.generateQueries(&ac, queries); err != nil {
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

func (r *RuleUDT) generateQueries(ruleResolvable *ResolvableUDT, queries *map[string]QueryUDT) error {
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
