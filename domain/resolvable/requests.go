package resolvable

import (
	"fmt"
	"ifttt/manager/common"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mitchellh/mapstructure"
)

const (
	accessorRuleResolvable                = "rule"
	accessorJqResolvable                  = "jq"
	accessorGetRequestResolvable          = "getReq"
	accessorGetResponseResolvable         = "getRes"
	accessorGetQueryResultsResolvable     = "getQueryRes"
	accessorGetApiResultsResolvable       = "getApiRes"
	accessorGetStoreResolvable            = "getStore"
	accessorGetConstResolvable            = "const"
	accessorArithmetic                    = "arithmetic"
	accessorQueryResolvable               = "query"
	accessorApiCallResolvable             = "api"
	accessorSetResResolvable              = "setRes"
	accessorSetStoreResolvable            = "setStore"
	accessorResponseResolvable            = "sendRes"
	accessorPreConfigResolvable           = "getPreConfig"
	accessorStringInterpolationResolvable = "stringInterpolation"
)

var resolveTypes = []string{
	accessorRuleResolvable,
	accessorJqResolvable,
	accessorGetRequestResolvable,
	accessorGetResponseResolvable,
	accessorGetQueryResultsResolvable,
	accessorGetApiResultsResolvable,
	accessorGetStoreResolvable,
	accessorGetConstResolvable,
	accessorArithmetic,
	accessorQueryResolvable,
	accessorApiCallResolvable,
	accessorSetResResolvable,
	accessorSetStoreResolvable,
	accessorResponseResolvable,
	accessorPreConfigResolvable,
	accessorStringInterpolationResolvable,
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
					return nil
				case accessorGetResponseResolvable:
					return nil
				case accessorGetQueryResultsResolvable:
					return nil
				case accessorGetApiResultsResolvable:
					return nil
				case accessorGetStoreResolvable:
					return nil
				case accessorGetConstResolvable:
					return nil
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
				case accessorResponseResolvable:
					resolver = &responseResolvable{}
				case accessorRuleResolvable:
					resolver = &callRuleResolvable{}
				case accessorPreConfigResolvable:
					resolver = &preConfigResolvable{}
				case accessorStringInterpolationResolvable:
					resolver = &stringInterpolationResolvable{}
				default:
					return validation.NewError("resolvable_not_found",
						fmt.Sprintf("resolvable %s not found", r.ResolveType))
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

func (c *callRuleResolvable) Validate() error {
	return validation.Validate(&c.RuleId, validation.Required)
}

func (c *apiCallResolvable) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Method, validation.Required, validation.In(
			http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut)),
		validation.Field(&c.Url, validation.Required, is.URL),
		validation.Field(&c.Headers, validation.By(func(value interface{}) error {
			m := value.(map[string]any)
			return validateMapIfResolvable(m)
		})),
		validation.Field(&c.Body, validation.By(func(value interface{}) error {
			m := value.(map[string]any)
			return validateMapIfResolvable(m)
		})),
	)
}

func (c *arithmetic) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Group, validation.NotNil),
		validation.Field(&c.Operation, validation.Required, validation.In("+", "-", "*", "/", "%")),
		validation.Field(&c.Operators, validation.Each(validation.Required, validation.By(
			func(value interface{}) error {
				a := value.(arithmetic)
				return a.Validate()
			}))),
		validation.Field(&c.Value, validation.When(!c.Group, validation.By(
			func(value interface{}) error {
				r := value.(Resolvable)
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
	namedParams := common.RegexNamedParameters.FindAllString(c.QueryHash, -1)
	namedCount := len(namedParams)
	positionalCount := len(common.RegexPositionalParameters.FindAllString(c.QueryHash, -1))
	return validation.ValidateStruct(c,
		validation.Field(&c.QueryString, validation.Required),
		validation.Field(&c.QueryHash, validation.Required),
		validation.Field(&c.Return, validation.Required),
		validation.Field(&c.Named, validation.Required),
		validation.Field(&c.NamedParameters, validation.When(c.Named, validation.Length(namedCount, namedCount),
			validation.By(func(value interface{}) error {
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

func (c *preConfigResolvable) Validate() error {
	var mapCasted map[string]Resolvable = *c
	for _, val := range mapCasted {
		if err := val.Validate(); err != nil {
			return err
		}
	}
	return nil
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
