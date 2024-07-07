package models

import (
	"errors"

	jsontocql "github.com/Hammad1029/json-to-cql"
	"github.com/mitchellh/mapstructure"
)

type Resolvable struct {
	ResolveType string                 `json:"resolveType" mapstructure:"resolveType"`
	ResolveData map[string]interface{} `json:"resolveData" mapstructure:"resolveData"`
}

func (ruleResolvable *Resolvable) generateQueries(queries *map[string]QueryUDT) error {
	if query, ok := ruleResolvable.ResolveData["query"]; ruleResolvable.ResolveType == "db" && ok {
		if queryMap, ok := query.(map[string]interface{}); ok {
			var queryDoc jsontocql.QueryDoc
			mapstructure.Decode(queryMap, &queryDoc)
			if paraQuery, err := queryDoc.CreateParameterizedQuery(); err != nil {
				return err
			} else {
				ruleResolvable.ResolveData["query"] = paraQuery.QueryHash

				var queryResolvables []Resolvable
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
