package resolvable

import (
	"fmt"
	"ifttt/manager/common"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
}

func (r *Resolvable) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ResolveType, validation.Required,
			validation.In(common.ConvertStringToInterfaceArray(resolveTypes)...)),
		validation.Field(&r.ResolveData, validation.By(
			func(value interface{}) error {
				var resolver common.ValidatorInterface
				switch r.ResolveType {
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
				default:
					return validation.NewError("resolvable_not_found",
						fmt.Sprintf("resolvable %s not found", r.ResolveType))
				}
				if resolver == nil {
					return nil
				}
				data := value.(map[string]any)
				if err := mapstructure.Decode(data, &resolver); err != nil {
					return validation.NewInternalError(err)
				}
				return resolver.Validate()
			}),
		),
	)
}

func (d *dbDumpResolvable) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Table, validation.Required, validation.Length(3, 0)),
		validation.Field(&d.Columns, validation.Required, validation.Each(
			validation.By(func(value interface{}) error {
				r := value.(Resolvable)
				return r.Validate()
			}))),
	)
}

func (u *uuidResolvable) Validate() error {
	return nil
}

func (g *getCacheResolvable) Validate() error {
	return validation.Validate(&g.Key, validation.By(func(value interface{}) error {
		r := value.(Resolvable)
		return r.Validate()
	}))
}

func (s *setCacheResolvable) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Key, validation.Required, validation.By(func(value interface{}) error {
			r := value.(Resolvable)
			return r.Validate()
		})),
		validation.Field(&s.Value, validation.Required, validation.By(func(value interface{}) error {
			r := value.(Resolvable)
			return r.Validate()
		})),
		validation.Field(&s.TTL),
	)
}

func (s *encodeResolvable) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Alg, validation.Required,
			validation.In(common.EncodeMD5, common.EncodeSHA1, common.EncodeSHA2,
				common.EncodeBcrypt, common.EncodeBase64Decode, common.EncodeBase64Encode),
		),
		validation.Field(&s.Input, validation.Required, validation.By(func(value interface{}) error {
			r := value.(Resolvable)
			return r.Validate()
		})))
}

func (s *setLogResolvable) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.LogType, validation.Required, validation.In(common.LogError, common.LogInfo)),
		validation.Field(&s.LogData, validation.NotNil, validation.By(
			func(value interface{}) error {
				if r, ok := value.(Resolvable); ok {
					return r.Validate()
				}
				return nil
			},
		)))
}

func (g *getRequestResolvable) Validate() error {
	return nil
}

func (g *getResponseResolvable) Validate() error {
	return nil
}

func (g *getStoreResolvable) Validate() error {
	return nil
}

func (g *getPreConfigResolvable) Validate() error {
	return nil
}

func (g *getHeadersResolvable) Validate() error {
	return nil
}

func (g *getConstResolvable) Validate() error {
	return nil
}

func (c *apiCallResolvable) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Method, validation.Required, validation.In(
			http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut)),
		validation.Field(&c.URL, validation.Required, validation.By(
			func(value interface{}) error {
				r := value.(Resolvable)
				return r.Validate()
			})),
		validation.Field(&c.Headers, validation.By(func(value interface{}) error {
			m := value.(map[string]any)
			return validateMapIfResolvable(m)
		})),
		validation.Field(&c.Body, validation.By(func(value interface{}) error {
			m := value.(map[string]any)
			return validateMapIfResolvable(m)
		})),
		validation.Field(&c.Aysnc),
		validation.Field(&c.Timeout),
	)
}

func (c *arithmetic) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Group, validation.NotNil),
		validation.Field(&c.Operation, validation.When(c.Group, validation.Required,
			validation.In("+", "-", "*", "/", "%")).Else(validation.Empty)),
		validation.Field(&c.Operators, validation.Each(validation.Required, validation.By(
			func(value interface{}) error {
				a := value.(arithmetic)
				return a.Validate()
			}))),
		validation.Field(&c.Value, validation.When(!c.Group, validation.By(
			func(value interface{}) error {
				r := value.(*Resolvable)
				return r.Validate()
			})).Else(validation.Nil)),
	)
}

func (c *jqResolvable) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Query, validation.Required, validation.By(
			func(value interface{}) error {
				r := value.(Resolvable)
				return r.Validate()
			})),
		validation.Field(&c.Input, validation.NotNil),
	)
}

func (c *stringInterpolationResolvable) Validate() error {
	parametersCount := len(common.RegexStringInterpolationParameters.FindAllString(c.Template, -1))
	return validation.ValidateStruct(c,
		validation.Field(&c.Template, validation.Required, validation.By(func(value interface{}) error {
			if parametersCount == 0 {
				return validation.NewError("template_string_wrong", "template string must have atleast 1 parameter")
			}
			return nil
		})),
		validation.Field(&c.Parameters, validation.Length(parametersCount, parametersCount),
			validation.Each(validation.Required, validation.By(func(value interface{}) error {
				r := value.(Resolvable)
				return r.Validate()
			}))),
	)
}

func (c *queryResolvable) Validate() error {
	namedParams := common.RegexNamedParameters.FindAllString(c.QueryString, -1)
	namedCount := len(namedParams)
	positionalCount := len(common.RegexPositionalParameters.FindAllString(c.QueryString, -1))
	return validation.ValidateStruct(c,
		validation.Field(&c.QueryString, validation.Required),
		validation.Field(&c.QueryHash, validation.Required, validation.In(common.GetMD5Hash(c.QueryString))),
		validation.Field(&c.Return, validation.NotNil),
		validation.Field(&c.Named, validation.NotNil),
		validation.Field(&c.NamedParameters, validation.When(c.Named,
			validation.Length(namedCount, namedCount), validation.By(
				func(value interface{}) error {
					paramMap := value.(map[string]Resolvable)
					for _, key := range namedParams {
						if r, ok := paramMap[key]; !ok {
							return validation.NewError("named_param_not_found",
								fmt.Sprintf("named parameter %s not found", key))
						} else if err := r.Validate(); err != nil {
							return err
						}
					}
					return nil
				})).Else(validation.Empty)),
		validation.Field(&c.PositionalParameters, validation.When(!c.Named,
			validation.Length(positionalCount, positionalCount), validation.Each(validation.By(
				func(value interface{}) error {
					param := value.(Resolvable)
					return param.Validate()
				},
			))).Else(validation.Empty)),
		validation.Field(&c.Async),
		validation.Field(&c.Timeout),
	)
}

func (c *responseResolvable) Validate() error {
	return nil
}

func (c *setResResolvable) Validate() error {
	var mapCasted map[string]any = *c
	return validateMapIfResolvable(mapCasted)
}

func (c *setStoreResolvable) Validate() error {
	var mapCasted map[string]any = *c
	return validateMapIfResolvable(mapCasted)

}

func validateMapIfResolvable(val map[string]any) error {
	for _, val := range val {
		var r Resolvable
		if err := mapstructure.Decode(val, &r); err != nil || r.ResolveType == "" {
			continue
		} else if err := r.Validate(); err != nil {
			return err
		}
	}
	return nil
}
