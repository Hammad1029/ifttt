package resolvable

import (
	"fmt"
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
)

const (
	accessorJq                  = "jq"
	accessorGetErrors           = "getErrors"
	accessorGetRequest          = "getReq"
	accessorGetResponse         = "getRes"
	accessorGetStore            = "getStore"
	accessorGetConst            = "const"
	accessorArithmetic          = "arithmetic"
	accessorQuery               = "query"
	accessorApiCall             = "api"
	accessorSetRes              = "setRes"
	accessorSetStore            = "setStore"
	accessorSetLog              = "log"
	accessorEvent               = "event"
	accessorPreConfig           = "getPreConfig"
	accessorStringInterpolation = "stringInterpolation"
	accessorEncode              = "encode"
	accessorSetCache            = "setCache"
	accessorGetCache            = "getCache"
	accessorUUID                = "uuid"
	accessorHeaders             = "headers"
	accessorCast                = "cast"
	accessorOrm                 = "orm"
	accessorForEach             = "forEach"
	accessorGetIter             = "getIter"
	accessorDateInput           = "dateInput"
	accessorDateManipulator     = "dateManipulator"
	accessorDateFunc            = "dateFunc"
)

var resolveTypes = []string{
	accessorJq,
	accessorGetErrors,
	accessorGetRequest,
	accessorGetResponse,
	accessorGetStore,
	accessorGetConst,
	accessorArithmetic,
	accessorQuery,
	accessorApiCall,
	accessorSetRes,
	accessorSetStore,
	accessorSetLog,
	accessorEvent,
	accessorPreConfig,
	accessorStringInterpolation,
	accessorEncode,
	accessorSetCache,
	accessorGetCache,
	accessorUUID,
	accessorHeaders,
	accessorCast,
	accessorOrm,
	accessorForEach,
	accessorGetIter,
	accessorDateInput,
	accessorDateManipulator,
	accessorDateFunc,
}

func factory(template any) (resolvableInterface, error) {
	var base Resolvable
	if err := mapstructure.Decode(template, &base); err != nil {
		return nil, err
	}

	var resolver resolvableInterface
	switch base.ResolveType {
	case accessorJq:
		resolver = &jq{}
	case accessorGetErrors:
		resolver = &getErrors{}
	case accessorGetRequest:
		resolver = &getRequest{}
	case accessorGetResponse:
		resolver = &getResponse{}
	case accessorGetStore:
		resolver = &getStore{}
	case accessorGetConst:
		resolver = &getConst{}
	case accessorArithmetic:
		resolver = &arithmetic{}
	case accessorQuery:
		resolver = &query{}
	case accessorApiCall:
		resolver = &apiCall{}
	case accessorSetRes:
		resolver = &setRes{}
	case accessorSetStore:
		resolver = &setStore{}
	case accessorSetLog:
		resolver = &setLog{}
	case accessorEvent:
		resolver = &event{}
	case accessorPreConfig:
		resolver = &getPreConfig{}
	case accessorStringInterpolation:
		resolver = &stringInterpolation{}
	case accessorEncode:
		resolver = &encode{}
	case accessorSetCache:
		resolver = &setCache{}
	case accessorGetCache:
		resolver = &getCache{}
	case accessorUUID:
		resolver = &uuid{}
	case accessorHeaders:
		resolver = &getHeaders{}
	case accessorCast:
		resolver = &cast{}
	case accessorOrm:
		resolver = &Orm{}
	case accessorForEach:
		resolver = &forEach{}
	case accessorGetIter:
		resolver = &getIter{}
	case accessorDateInput:
		resolver = &dateInput{}
	case accessorDateManipulator:
		resolver = &dateManipulator{}
	case accessorDateFunc:
		resolver = &dateFunc{}
	default:
		return nil, fmt.Errorf("resolvable %s not found", base.ResolveType)
	}

	if err := mapstructure.Decode(base.ResolveData, &resolver); err != nil {
		return nil, fmt.Errorf("could not decode resolver %s: %s", base.ResolveType, err)
	}

	return resolver, nil
}

func ManipulateIfResolvable(val any, dependencies map[common.IntIota]any) (any, error) {
	if val == nil {
		return nil, nil
	}

	concrete := reflect.Indirect(reflect.ValueOf(val)).Interface()
	if concrete == nil {
		return nil, nil
	}

	switch reflect.TypeOf(concrete).Kind() {
	case reflect.Struct:
		{
			if r, ok := concrete.(common.Manipulatable); ok {
				if err := r.Manipulate(dependencies); err != nil {
					return nil, err
				} else {
					return r, nil
				}
			}
		}
	case reflect.Map:
		{
			var nested Resolvable
			err := mapstructure.Decode(concrete, &nested)
			if err == nil && nested.ResolveType != "" && nested.ResolveData != nil {
				if err := nested.Manipulate(dependencies); err != nil {
					return nil, err
				} else {
					return nested, nil
				}
			}

			var mapCloned map[string]any
			if err := mapstructure.Decode(concrete, &mapCloned); err != nil {
				return nil, err
			}
			for key := range mapCloned {
				val := mapCloned[key]
				if v, err := ManipulateIfResolvable(&val, dependencies); err != nil {
					return nil, err
				} else {
					mapCloned[key] = v
				}
			}
			return mapCloned, nil
		}
	case reflect.Slice, reflect.Array:
		{
			var oArr []any
			if err := mapstructure.Decode(concrete, &oArr); err != nil {
				return nil, err
			}
			for idx, item := range oArr {
				if v, err := ManipulateIfResolvable(&item, dependencies); err != nil {
					return nil, err
				} else {
					oArr[idx] = v
				}
			}
			return oArr, nil

		}
	}
	return concrete, nil
}

func ValidateIfResolvable(val any) error {
	if val == nil {
		return nil
	}

	concrete := reflect.Indirect(reflect.ValueOf(val)).Interface()
	if concrete == nil {
		return nil
	}

	switch reflect.TypeOf(concrete).Kind() {
	case reflect.Struct:
		{
			if r, ok := concrete.(common.Validatable); ok {
				return r.Validate()
			}
		}
	case reflect.Map:
		{
			var nested Resolvable
			err := mapstructure.Decode(concrete, &nested)
			if err == nil && nested.ResolveType != "" && nested.ResolveData != nil {
				return nested.Validate()
			}

			var mapCloned map[string]any
			if err := mapstructure.Decode(concrete, &mapCloned); err != nil {
				return err
			}
			for key := range mapCloned {
				val := mapCloned[key]
				if err = ValidateIfResolvable(&val); err != nil {
					return err
				}
			}
		}
	case reflect.Slice, reflect.Array:
		{
			var oArr []any
			if err := mapstructure.Decode(concrete, &oArr); err != nil {
				return err
			}
			for _, item := range oArr {
				if err := ValidateIfResolvable(&item); err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func mustBeResolvable(val any) error {
	var r Resolvable
	if err := mapstructure.Decode(val, &r); err != nil || r.ResolveType == "" {
		return fmt.Errorf("provided object is not resolvable")
	} else if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

func checkIfResolvable(val any) *Resolvable {
	var r Resolvable
	if err := mapstructure.Decode(val, &r); err != nil || r.ResolveType == "" {
		return nil
	}
	return &r
}

func ManipulateArray(arr []Resolvable, dependencies map[common.IntIota]any) ([]Resolvable, error) {
	var manipulated []Resolvable
	for _, r := range arr {
		if err := r.Manipulate(dependencies); err != nil {
			return nil, err
		}
		manipulated = append(manipulated, r)
	}
	return manipulated, nil
}

func ManipulateMap(arr map[string]Resolvable, dependencies map[common.IntIota]any) (map[string]Resolvable, error) {
	manipulated := make(map[string]Resolvable)
	for key, r := range arr {
		if err := r.Manipulate(dependencies); err != nil {
			return nil, err
		}
		manipulated[key] = r
	}
	return manipulated, nil
}

func (p *Orm) ManipulatePopulate(
	populate *[]orm_schema.Populate, models *map[string]*orm_schema.Model, dependencies map[common.IntIota]any,
) error {
	for _, child := range *populate {
		if pModel, ok := (*models)[p.Model]; !ok {
			return fmt.Errorf("model in populate not found: %s", p.Model)
		} else if pModel.PrimaryKey == "" {
			return fmt.Errorf("model %s does not contain primary key for populate", p.Model)
		} else if err := p.ManipulateProjection(child.Model, models, &child.Project); err != nil {
			return err
		} else if err := p.ManipulateWhere(&child.Where, dependencies); err != nil {
			return err
		}

		if err := p.ManipulatePopulate(&child.Populate, models, dependencies); err != nil {
			return err
		}
	}
	return nil
}

func (w *Orm) ManipulateWhere(where *orm_schema.Where, dependencies map[common.IntIota]any) error {
	if manipulated, err := ManipulateIfResolvable(where.Values, dependencies); err != nil {
		return err
	} else if mapped, ok := manipulated.([]any); !ok {
		return fmt.Errorf("could not cast conditionsValues to map")
	} else {
		where.Values = mapped
	}

	for _, v := range where.Values {
		if converted, err := anyToResolvable(v); err != nil {
			return err
		} else {
			w.Query.PositionalParameters = append(w.Query.PositionalParameters, *converted)
		}
	}
	return nil
}

func anyToResolvable(v any) (*Resolvable, error) {
	if res := checkIfResolvable(v); res != nil {
		return res, nil
	} else {
		constRes := getConst{Value: v}
		param := Resolvable{ResolveType: accessorGetConst}
		if err := mapstructure.Decode(constRes, &param.ResolveData); err != nil {
			return nil, err
		}
		return &param, nil
	}
}

func (o *Orm) ManipulateProjection(
	modelName string, models *map[string]*orm_schema.Model, projections *[]orm_schema.Projection,
) error {
	model, ok := (*models)[modelName]
	if !ok {
		return fmt.Errorf("model %s not found", modelName)
	}

	modelProjectionsMapped := lo.SliceToMap(model.Projections,
		func(p orm_schema.Projection) (string, orm_schema.Projection) {
			return p.Column, p
		})
	modelProjectionsKeys := lo.Keys(modelProjectionsMapped)
	customProjectionsMapped := lo.SliceToMap(*projections,
		func(p orm_schema.Projection) (string, orm_schema.Projection) {
			return p.Column, p
		})
	customProjectionsKeys := lo.Keys(customProjectionsMapped)

	if len(lo.Intersect(modelProjectionsKeys, customProjectionsKeys)) != len(customProjectionsKeys) {
		return fmt.Errorf("invalid projections")
	}

	for idx, p := range *projections {
		(*projections)[idx].DataType = modelProjectionsMapped[p.Column].DataType
	}

	return nil
}
