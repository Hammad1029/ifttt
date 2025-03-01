package resolvable

import (
	"fmt"
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/samber/lo"
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

func (r *apiCall) Manipulate(dependencies map[common.IntIota]any) error {
	if err := r.URL.Manipulate(dependencies); err != nil {
		return err
	}

	if manipulated, err := ManipulateIfResolvable(r.Headers, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast headers to map")
	} else {
		r.Headers = mapped
	}

	if manipulated, err := ManipulateIfResolvable(r.Body, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast body to map")
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

func (r *setCache) Manipulate(dependencies map[common.IntIota]any) error {
	if err := r.Key.Manipulate(dependencies); err != nil {
		return err
	} else if err := r.Value.Manipulate(dependencies); err != nil {
		return err
	}
	return nil
}

func (r *getCache) Manipulate(dependencies map[common.IntIota]any) error {
	return r.Key.Manipulate(dependencies)
}

func (r *deleteCache) Manipulate(dependencies map[common.IntIota]any) error {
	return r.Key.Manipulate(dependencies)
}

func (r *encode) Manipulate(dependencies map[common.IntIota]any) error {
	return r.Input.Manipulate(dependencies)
}

func (r *getErrors) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *getStore) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateIfResolvable(r.Query, dependencies); err != nil {
		return err
	} else {
		r.Query = manipulated
	}

	return nil
}

func (r *getHeaders) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *getConst) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *jq) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateIfResolvable(&r.Input, dependencies); err != nil {
		return err
	} else {
		r.Input = manipulated
	}

	if manipulated, err := ManipulateIfResolvable(r.Query, dependencies); err != nil {
		return err
	} else {
		r.Query = manipulated
	}

	return nil
}

func (r *query) Manipulate(dependencies map[common.IntIota]any) error {
	if _, err := sqlparser.Parse(r.QueryString); err != nil {
		return fmt.Errorf("parsing query failed: %s", err)
	}

	if manipulated, err := ManipulateArray(&r.Parameters, dependencies); err != nil {
		return err
	} else {
		r.Parameters = *manipulated
	}

	return nil
}

func (r *response) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *setStore) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateIfResolvable(r, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.(map[string]any); !ok {
		return fmt.Errorf("could not cast setstore to map")
	} else {
		*r = setStore(mapped)
	}
	return nil
}

func (r *setLog) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *stringInterpolation) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateArray(&r.Parameters, dependencies); err != nil {
		return err
	} else {
		r.Parameters = *manipulated
		return nil
	}
}

func (r *uuid) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (r *cast) Manipulate(dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateIfResolvable(r.Input, dependencies); err != nil {
		return err
	} else {
		r.Input = manipulated
	}
	return nil
}

func (r *Orm) Manipulate(dependencies map[common.IntIota]any) error {
	switch r.Operation {
	case common.OrmSelect, common.OrmInsert, common.OrmUpdate, common.OrmDelete:
	default:
		return fmt.Errorf("operation %s not allowed", r.Operation)
	}

	queryGenerator, ok := dependencies[common.DependencyOrmQueryRepo].(OrmQueryGenerator)
	if !ok {
		return fmt.Errorf("could not get query generator repo")
	}
	ormRepo, ok := dependencies[common.DependencyOrmSchemaRepo].(orm_schema.OrmRepository)
	if !ok {
		return fmt.Errorf("could not get orm repo")
	}

	r.Query = &query{
		Parameters: []Resolvable{},
	}

	allModels, err := ormRepo.GetAllModels()
	if err != nil {
		return err
	}
	rootModel, ok := allModels[r.Model]
	if !ok {
		return fmt.Errorf("model %s not found", r.Model)
	}

	var colSq []string
	if r.Columns != nil {
		if sq, err := r.ManipulateColumns(rootModel, dependencies); err != nil {
			return err
		} else {
			colSq = sq
		}
	}

	if r.Project != nil {
		if err := r.ManipulateProjection(r.Model, &allModels, r.Project); err != nil {
			return err
		}
	} else {
		r.Project = &[]orm_schema.Projection{}
	}

	if r.Where != nil {
		if err := r.ManipulateWhere(r.Where, dependencies); err != nil {
			return err
		}
	} else {
		r.Where = &orm_schema.Where{}
	}

	r.ModelsInUse = &[]string{}
	*r.ModelsInUse = append(*r.ModelsInUse, r.Model)
	if r.Populate != nil {
		if err := r.ManipulatePopulate(r.Populate, &allModels, dependencies); err != nil {
			return err
		}
	} else {
		r.Populate = &[]orm_schema.Populate{}
	}
	*r.ModelsInUse = lo.Uniq(*r.ModelsInUse)

	switch r.Operation {
	case common.OrmSelect:
		r.Query.Scan = true
		r.Query.QueryString, err = queryGenerator.GenerateSelect(r, rootModel, allModels)
	case common.OrmInsert:
		r.Query.QueryString, err = queryGenerator.GenerateInsert(rootModel.Table, colSq)
	case common.OrmUpdate:
		r.Query.QueryString, err = queryGenerator.GenerateUpdate(rootModel.Table, r.Where.Template, colSq)
	case common.OrmDelete:
		r.Query.QueryString, err = queryGenerator.GenerateDelete(rootModel.Table, r.Where.Template)
	default:
		return fmt.Errorf("generator for %s not available", r.Operation)
	}

	if err := r.Query.Manipulate(dependencies); err != nil {
		return err
	}

	if r.Operation == common.OrmInsert || r.Operation == common.OrmUpdate {
		r.SuccessiveQuery = &query{
			Scan: true,
		}
		if queryString, params, err := queryGenerator.GenerateSuccessive(r, rootModel); err != nil {
			return fmt.Errorf("could not generate successive query: %s", err)
		} else {
			r.SuccessiveQuery.QueryString = queryString
			if params != nil && len(*params) > 0 {
				r.SuccessiveQuery.Parameters = *params
			}
		}

		if err := r.SuccessiveQuery.Manipulate(dependencies); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (d *dateFunc) Manipulate(dependencies map[common.IntIota]any) error {
	if err := d.Input.Manipulate(dependencies); err != nil {
		return err
	}
	for _, m := range d.Manipulators {
		if err := m.Manipulate(dependencies); err != nil {
			return err
		}
	}
	return nil
}

func (d *dateManipulator) Manipulate(dependencies map[common.IntIota]any) error {
	return d.Operand.Manipulate(dependencies)
}

func (d *dateInput) Manipulate(dependencies map[common.IntIota]any) error {
	if d.Input != nil {
		return d.Input.Manipulate(dependencies)
	}
	return nil
}

func (f *filterMap) Manipulate(dependencies map[common.IntIota]any) error {
	if err := f.Condition.Manipulate(dependencies); err != nil {
		return err
	}

	if manipulated, err := ManipulateIfResolvable(f.Input, dependencies); err != nil {
		return err
	} else {
		f.Input = manipulated
	}

	if manipulated, err := ManipulateArray(f.Do, dependencies); err != nil {
		return err
	} else {
		f.Do = manipulated
	}

	return nil
}

func (f *getIter) Manipulate(dependencies map[common.IntIota]any) error {
	return nil
}

func (c *Condition) Manipulate(dependencies map[common.IntIota]any) error {
	if c.Group {
		for _, cnd := range c.Conditions {
			if err := cnd.Manipulate(dependencies); err != nil {
				return err
			}
		}
	} else if err := c.Operator1.Manipulate(dependencies); err != nil {
		return err
	} else if err := c.Operator2.Manipulate(dependencies); err != nil {
		return err
	}
	return nil
}

func (c *conditional) Manipulate(dependencies map[common.IntIota]any) error {
	if err := c.Condition.Manipulate(dependencies); err != nil {
		return err
	}

	if manipulated, err := ManipulateArray(&c.True, dependencies); err != nil {
		return err
	} else {
		c.True = *manipulated
	}

	if manipulated, err := ManipulateArray(&c.False, dependencies); err != nil {
		return err
	} else {
		c.False = *manipulated
	}

	return nil
}

func (d *dateIntervals) Manipulate(dependencies map[common.IntIota]any) error {
	if err := d.Start.Manipulate(dependencies); err != nil {
		return err
	} else if err := d.End.Manipulate(dependencies); err != nil {
		return err
	}

	return nil
}
