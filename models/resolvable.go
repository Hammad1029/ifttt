package models

import (
	"errors"

	jsontocql "github.com/Hammad1029/json-to-cql"
	"github.com/mitchellh/mapstructure"
)

type ResolvableUDT struct {
	Type string                 `cql:"type" json:"type" mapstructure:"type"`
	Data map[string]interface{} `cql:"data" json:"data" mapstructure:"data"`
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
