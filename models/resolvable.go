package models

import (
	"errors"
	"fmt"

	jsontocql "github.com/Hammad1029/json-to-cql"
	"github.com/mitchellh/mapstructure"
)

type Resolvable struct {
	ResolveType string                 `json:"resolveType" mapstructure:"resolveType"`
	ResolveData map[string]interface{} `json:"resolveData" mapstructure:"resolveData"`
}

func (r *Resolvable) transformResolvables(queries *map[string]jsontocql.ParameterizedQuery) error {
	if r.ResolveType == "db" {
		err := r.generateQueries(queries)
		return err
	}

	for key, v := range r.ResolveData {
		switch value := v.(type) {
		case map[string]interface{}:
			nestedResolvable := Resolvable{}
			if err := mapstructure.Decode(value, &nestedResolvable); err == nil {
				err := nestedResolvable.transformResolvables(queries)
				r.ResolveData[key] = nestedResolvable
				return err
			}
		}
	}
	return nil
}

func (r *Resolvable) generateQueries(queries *map[string]jsontocql.ParameterizedQuery) error {
	if query, ok := r.ResolveData["query"]; ok {
		var queryDoc jsontocql.QueryDoc
		if err := mapstructure.Decode(query, &queryDoc); err != nil {
			return fmt.Errorf("method generateQueries: couldn't decode query doc | %s", err.Error())
		}
		if paraQuery, err := queryDoc.CreateParameterizedQuery(); err != nil {
			return fmt.Errorf("method generateQueries: couldn't create parameterized query | %s", err.Error())
		} else {
			r.ResolveData["query"] = paraQuery.QueryHash

			(*queries)[paraQuery.QueryHash] = jsontocql.ParameterizedQuery{
				QueryString: paraQuery.QueryString,
				QueryHash:   paraQuery.QueryHash,
				Resolvables: paraQuery.Resolvables,
				Type:        paraQuery.Type,
			}
			return nil
		}
	} else {
		return errors.New("method generateQueries: no query found in resolvableType db")
	}
}
