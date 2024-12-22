package resolvable

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

const (
	accessorJqResolvable                  = "jq"
	accessorGetRequestResolvable          = "getReq"
	accessorGetResponseResolvable         = "getRes"
	accessorGetStoreResolvable            = "getStore"
	accessorGetConstResolvable            = "const"
	accessorArithmetic                    = "arithmetic"
	accessorQueryResolvable               = "query"
	accessorApiCallResolvable             = "api"
	accessorSetResResolvable              = "setRes"
	accessorSetStoreResolvable            = "setStore"
	accessorSetLogResolvable              = "log"
	accessorResponseResolvable            = "sendRes"
	accessorPreConfigResolvable           = "getPreConfig"
	accessorStringInterpolationResolvable = "stringInterpolation"
	accessorEncodeResolvable              = "encode"
	accessorSetCacheResolvable            = "setCache"
	accessorGetCacheResolvable            = "getCache"
	accessorUUIDResolvable                = "uuid"
	accessorHeadersResolvable             = "headers"
	accessorDbDumpResolvable              = "dbDump"
	accessorCastResolvable                = "cast"
	accessorOrmResolvable                 = "orm"
)

var resolveTypes = []string{
	accessorJqResolvable,
	accessorGetRequestResolvable,
	accessorGetResponseResolvable,
	accessorGetStoreResolvable,
	accessorGetConstResolvable,
	accessorArithmetic,
	accessorQueryResolvable,
	accessorApiCallResolvable,
	accessorSetResResolvable,
	accessorSetStoreResolvable,
	accessorSetLogResolvable,
	accessorResponseResolvable,
	accessorPreConfigResolvable,
	accessorStringInterpolationResolvable,
	accessorEncodeResolvable,
	accessorSetCacheResolvable,
	accessorGetCacheResolvable,
	accessorUUIDResolvable,
	accessorHeadersResolvable,
	accessorDbDumpResolvable,
	accessorCastResolvable,
	accessorOrmResolvable,
}

func factory(template any) (ResolvableInterface, error) {
	var base Resolvable
	if err := mapstructure.Decode(template, &base); err != nil {
		return nil, err
	}

	var resolver ResolvableInterface
	switch base.ResolveType {
	case accessorJqResolvable:
		resolver = &jqResolvable{}
	case accessorGetRequestResolvable:
		resolver = &getRequestResolvable{}
	case accessorGetResponseResolvable:
		resolver = &getResponseResolvable{}
	case accessorGetStoreResolvable:
		resolver = &getStoreResolvable{}
	case accessorGetConstResolvable:
		resolver = &getConstResolvable{}
	case accessorArithmetic:
		resolver = &arithmetic{}
	case accessorQueryResolvable:
		resolver = &queryResolvable{}
	case accessorApiCallResolvable:
		resolver = &apiCallResolvable{}
	case accessorSetResResolvable:
		resolver = &setResResolvable{}
	case accessorSetStoreResolvable:
		resolver = &setStoreResolvable{}
	case accessorSetLogResolvable:
		resolver = &setLogResolvable{}
	case accessorResponseResolvable:
		resolver = &responseResolvable{}
	case accessorPreConfigResolvable:
		resolver = &getPreConfigResolvable{}
	case accessorStringInterpolationResolvable:
		resolver = &stringInterpolationResolvable{}
	case accessorEncodeResolvable:
		resolver = &encodeResolvable{}
	case accessorSetCacheResolvable:
		resolver = &setCacheResolvable{}
	case accessorGetCacheResolvable:
		resolver = &getCacheResolvable{}
	case accessorUUIDResolvable:
		resolver = &uuidResolvable{}
	case accessorHeadersResolvable:
		resolver = &getHeadersResolvable{}
	case accessorDbDumpResolvable:
		resolver = &dbDumpResolvable{}
	case accessorCastResolvable:
		resolver = &castResolvable{}
	case accessorOrmResolvable:
		resolver = &ormResolvable{}
	default:
		return nil, fmt.Errorf("resolvable %s not found", base.ResolveType)
	}

	if err := mapstructure.Decode(base.ResolveData, &resolver); err != nil {
		return nil, err
	}

	return resolver, nil
}

func ifNestedResolvable(val any, validate bool) error {
	var err error

	switch o := val.(type) {
	case nil:
		return nil
	case ResolvableInterface:
		if validate {
			return o.Validate()
		} else {
			return o.Manipulate()
		}
	default:
		{
			switch reflect.TypeOf(o).Kind() {
			case reflect.Map:
				{
					var nested Resolvable
					err = mapstructure.Decode(o, &nested)
					if err == nil && nested.ResolveType != "" && nested.ResolveData != nil {
						if validate {
							return nested.Validate()
						} else {
							return nested.Manipulate()
						}
					}

					var mapCloned map[string]any
					if err := mapstructure.Decode(o, &mapCloned); err != nil {
						return err
					}
					for _, val := range mapCloned {
						if err = ifNestedResolvable(val, validate); err != nil {
							return err
						}
					}
					return nil
				}
			case reflect.Slice, reflect.Array:
				{
					var oArr []any
					if err := mapstructure.Decode(o, &oArr); err != nil {
						return err
					}
					for _, item := range oArr {
						if err = ifNestedResolvable(item, validate); err != nil {
							return err
						}
					}
					return nil

				}
			default:
				return nil
			}
		}
	}
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

func ManipulateArray(arr []Resolvable) error {
	for _, r := range arr {
		if err := r.Manipulate(); err != nil {
			return err
		}
	}
	return nil
}
