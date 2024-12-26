package resolvable

import (
	"fmt"
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"

	"github.com/go-viper/mapstructure/v2"
)

func (r *Resolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if resolver, err := factory(r); err != nil {
		return err
	} else if err := resolver.Manipulate(dependencies); err != nil {
		return err
	} else if mapped, err := common.AnyToMap(resolver); err != nil {
		return fmt.Errorf("could not convert to map: %s", err)
	} else {
		r.ResolveData = mapped
		return nil
	}
}

func (r *apiCallResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if err := r.URL.Manipulate(dependencies); err != nil {
		return err
	}

	if manipulated, err := manipulateIfResolvable(r.Headers, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast headers to map")
	} else {
		r.Headers = mapped
	}

	if manipulated, err := manipulateIfResolvable(r.Body, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast to map")
	} else {
		r.Headers = mapped
	}

	return nil
}

func (r *arithmetic) Manipulate(dependencies map[common.IntIota]any) error {
	if r.Group {
		for _, a := range r.Operators {
			if err := a.Manipulate(dependencies); err != nil {
				return err
			}
		}
	} else {
		return r.Value.Manipulate(dependencies)
	}
	return nil
}

func (r *setCacheResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if err := r.Key.Manipulate(dependencies); err != nil {
		return err
	} else if err := r.Value.Manipulate(dependencies); err != nil {
		return err
	}
	return nil
}

func (r *getCacheResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return r.Key.Manipulate(dependencies)
}

func (r *dbDumpResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateMap(r.Columns, dependencies); err != nil {
		return err
	} else {
		r.Columns = manipulated
		return nil
	}
}

func (r *encodeResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return r.Input.Manipulate(dependencies)
}

func (r *getRequestResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *getResponseResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *getStoreResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *getPreConfigResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *getHeadersResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *getConstResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *jqResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := manipulateIfResolvable(&r.Input, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast to map")
	} else {
		r.Input = mapped
	}

	if manipulated, err := manipulateIfResolvable(r.Query, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast to map")
	} else {
		r.Query = mapped
	}

	return nil
}

func (r *queryResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateMap(r.NamedParameters, dependencies); err != nil {
		return err
	} else {
		r.NamedParameters = manipulated
	}

	if manipulated, err := ManipulateArray(r.PositionalParameters, dependencies); err != nil {
		return err
	} else {
		r.PositionalParameters = manipulated
	}

	return nil
}

func (r *responseResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *setResResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := manipulateIfResolvable(r, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast to map")
	} else {
		*r = setResResolvable(mapped)
	}
	return nil
}

func (r *setStoreResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := manipulateIfResolvable(r, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast to map")
	} else {
		*r = setStoreResolvable(mapped)
	}
	return nil
}

func (r *setLogResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *stringInterpolationResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateArray(r.Parameters, dependencies); err != nil {
		return err
	} else {
		r.Parameters = manipulated
		return nil
	}
}

func (r *uuidResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *castResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := manipulateIfResolvable(r.Input, dependencies); err != nil {
		return err
	} else {
		r.Input = manipulated
	}
	return nil
}

func (r *OrmResolvable) Manipulate(dependencies map[common.IntIota]any) error {
	queryGenerator, ok := dependencies[common.DependencyOrmQueryRepo].(OrmQueryGenerator)
	if !ok {
		return fmt.Errorf("could not get query generator repo")
	}
	ormRepo, ok := dependencies[common.DependencyOrmSchemaRepo].(orm_schema.OrmRepository)
	if !ok {
		return fmt.Errorf("could not get orm repo")
	}

	if manipulated, err := manipulateIfResolvable(r.ConditionsValue, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.([]any); !ok {
		return fmt.Errorf("could not cast to map")
	} else {
		r.ConditionsValue = mapped
	}

	r.Query = &queryResolvable{
		NamedParameters:      map[string]Resolvable{},
		PositionalParameters: []Resolvable{},
	}
	if r.Operation == common.OrmSelect {
		r.Query.Return = true
	} else {
		return fmt.Errorf("operation %s not allowed", r.Operation)
	}

	allModels, err := ormRepo.GetAllModels()
	if err != nil {
		return err
	}

	if queryString, err := queryGenerator.Generate(r, allModels); err != nil {
		return err
	} else {
		r.Query.QueryString = queryString
	}

	for _, v := range r.ConditionsValue {
		if res := checkIfResolvable(v); res != nil {
			r.Query.PositionalParameters = append(r.Query.PositionalParameters, *res)
		} else {
			constRes := getConstResolvable{Value: v}
			param := Resolvable{ResolveType: accessorGetConstResolvable}
			if err := mapstructure.Decode(constRes, &param.ResolveData); err != nil {
				return err
			}
			r.Query.PositionalParameters = append(r.Query.PositionalParameters, param)
		}
	}

	return nil
}
