package resolvable

import (
	"ifttt/manager/common"

	"github.com/go-viper/mapstructure/v2"
)

func (r *Resolvable) Manipulate() error {
	if resolver, err := factory(r); err != nil {
		return err
	} else {
		return resolver.Manipulate()
	}
}

func (r *apiCallResolvable) Manipulate() error {
	if err := r.URL.Manipulate(); err != nil {
		return err
	}
	if err := ifNestedResolvable(r.Headers, false); err != nil {
		return err
	} else if err := ifNestedResolvable(r.Body, false); err != nil {
		return err
	}
	return nil
}

func (r *arithmetic) Manipulate() error {
	if r.Group {
		for _, a := range r.Operators {
			if err := a.Manipulate(); err != nil {
				return err
			}
		}
	} else {
		return r.Value.Manipulate()
	}
	return nil
}

func (r *setCacheResolvable) Manipulate() error {
	if err := r.Key.Manipulate(); err != nil {
		return err
	} else if err := r.Value.Manipulate(); err != nil {
		return err
	}
	return nil
}

func (r *getCacheResolvable) Manipulate() error {
	return r.Key.Manipulate()
}

func (r *dbDumpResolvable) Manipulate() error {
	return ifNestedResolvable(r.Columns, false)
}

func (r *encodeResolvable) Manipulate() error {
	return r.Input.Manipulate()
}

func (r *getRequestResolvable) Manipulate() error {
	return nil
}

func (r *getResponseResolvable) Manipulate() error {
	return nil
}

func (r *getStoreResolvable) Manipulate() error {
	return nil
}

func (r *getPreConfigResolvable) Manipulate() error {
	return nil
}

func (r *getHeadersResolvable) Manipulate() error {
	return nil
}

func (r *getConstResolvable) Manipulate() error {
	return nil
}

func (r *jqResolvable) Manipulate() error {
	if err := ifNestedResolvable(r.Input, false); err != nil {
		return err
	} else if err := ifNestedResolvable(r.Query, false); err != nil {
		return err
	}
	return nil
}

func (r *queryResolvable) Manipulate() error {
	if err := ifNestedResolvable(r.NamedParameters, false); err != nil {
		return err
	} else if err := ifNestedResolvable(r.PositionalParameters, false); err != nil {
		return err
	}
	return nil
}

func (r *responseResolvable) Manipulate() error {
	return nil
}

func (r *setResResolvable) Manipulate() error {
	return ifNestedResolvable(r, false)
}

func (r *setStoreResolvable) Manipulate() error {
	return ifNestedResolvable(r, false)
}

func (r *setLogResolvable) Manipulate() error {
	return nil
}

func (r *stringInterpolationResolvable) Manipulate() error {
	return ifNestedResolvable(r.Parameters, false)
}

func (r *uuidResolvable) Manipulate() error {
	return nil
}

func (r *castResolvable) Manipulate() error {
	return ifNestedResolvable(r.Input, false)
}

func (r *ormResolvable) Manipulate() error {
	if err := ifNestedResolvable(r.ConditionsValue, false); err != nil {
		return err
	}

	r.Query = queryResolvable{}
	if r.Operation == common.OrmSelect {
		r.Query.Return = true
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
