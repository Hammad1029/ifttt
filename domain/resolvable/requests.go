package resolvable

import (
	"fmt"
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r *Resolvable) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ResolveType, validation.Required,
			validation.In(common.ConvertStringToInterfaceArray(resolveTypes)...)),
		validation.Field(&r.ResolveData, validation.By(
			func(value any) error {
				if resolver, err := factory(r); err != nil {
					return validation.NewError("invalid_resolvabke", err.Error())
				} else {
					return resolver.Validate()
				}
			}),
		),
	)
}

func (d *dbDumpResolvable) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Table, validation.Required, validation.Length(3, 0)),
		validation.Field(&d.Columns, validation.Required, validation.Each(
			validation.By(func(value any) error {
				r, ok := value.(Resolvable)
				if !ok {
					return validation.NewError("resolvable_not_casted", "could not cast resolvable")
				}

				return r.Validate()
			}))),
	)
}

func (u *uuidResolvable) Validate() error {
	return nil
}

func (g *getCacheResolvable) Validate() error {
	return validation.Validate(&g.Key, validation.By(func(value any) error {
		return mustBeResolvable(value)
	}))
}

func (s *setCacheResolvable) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Key, validation.Required, validation.By(func(value any) error {
			r, ok := value.(Resolvable)
			if !ok {
				return validation.NewError("resolvable_not_casted", "could not cast resolvable")
			}

			return r.Validate()
		})),
		validation.Field(&s.Value, validation.Required, validation.By(func(value any) error {
			r, ok := value.(Resolvable)
			if !ok {
				return validation.NewError("resolvable_not_casted", "could not cast resolvable")
			}

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
		validation.Field(&s.Input, validation.Required, validation.By(func(value any) error {
			r, ok := value.(Resolvable)
			if !ok {
				return validation.NewError("resolvable_not_casted", "could not cast resolvable")
			}

			return r.Validate()
		})))
}

func (s *setLogResolvable) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.LogType, validation.Required, validation.In(common.LogError, common.LogInfo)),
		validation.Field(&s.LogData, validation.NotNil, validation.By(
			func(value any) error {
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
			func(value any) error {
				r, ok := value.(Resolvable)
				if !ok {
					return validation.NewError("resolvable_not_casted", "could not cast resolvable")
				}

				return r.Validate()
			})),
		validation.Field(&c.Headers, validation.By(func(value any) error {
			m := value.(map[string]any)
			return validateIfResolvable(m, nil)
		})),
		validation.Field(&c.Body, validation.By(func(value any) error {
			m := value.(map[string]any)
			return validateIfResolvable(m, nil)
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
			func(value any) error {
				a := value.(arithmetic)
				return a.Validate()
			}))),
		validation.Field(&c.Value, validation.When(!c.Group, validation.By(
			func(value any) error {
				r, ok := value.(*Resolvable)
				if !ok {
					return validation.NewError("resolvable_not_casted", "could not cast resolvable")
				}

				return r.Validate()
			})).Else(validation.Nil)),
	)
}

func (c *jqResolvable) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Query, validation.Required, validation.By(
			func(value any) error {
				if _, ok := value.(string); ok {
					return nil
				}
				if r, ok := value.(map[string]any); ok {
					return mustBeResolvable(r)
				}
				return validation.NewError("jq_invalid_query", "jq query must be string or resolvable")
			})),
		validation.Field(&c.Input, validation.NotNil),
	)
}

func (c *stringInterpolationResolvable) Validate() error {
	parametersCount := len(common.RegexStringInterpolationParameters.FindAllString(c.Template, -1))
	return validation.ValidateStruct(c,
		validation.Field(&c.Template, validation.Required, validation.By(func(value any) error {
			if parametersCount == 0 {
				return validation.NewError("template_string_wrong", "template string must have atleast 1 parameter")
			}
			return nil
		})),
		validation.Field(&c.Parameters, validation.Length(parametersCount, parametersCount),
			validation.Each(validation.Required, validation.By(func(value any) error {
				r, ok := value.(Resolvable)
				if !ok {
					return validation.NewError("resolvable_not_casted", "could not cast resolvable")
				}

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
		validation.Field(&c.Return, validation.NotNil),
		validation.Field(&c.Named, validation.NotNil),
		validation.Field(&c.NamedParameters, validation.When(c.Named,
			validation.Length(namedCount, namedCount), validation.By(
				func(value any) error {
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
				func(value any) error {
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
	return validateIfResolvable(mapCasted, nil)
}

func (c *setStoreResolvable) Validate() error {
	var mapCasted map[string]any = *c
	return validateIfResolvable(mapCasted, nil)
}

func (c *castResolvable) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.To, validation.In(common.ConvertStringToInterfaceArray(
			[]string{common.CastToString, common.CastToNumber, common.CastToBoolean})...)),
		validation.Field(&c.Input, validation.By(func(value any) error {
			return validateIfResolvable(value, nil)
		})),
	)
}

func (o *OrmResolvable) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Query, validation.Nil),
		validation.Field(&o.Operation, validation.In(
			common.OrmSelect, common.OrmUpdate, common.OrmInsert, common.OrmDelete)),
		validation.Field(&o.Model, validation.NotNil),
		validation.Field(&o.ConditionsTemplate, validation.NotNil),
		validation.Field(&o.ConditionsValue, validation.Each(validation.By(
			func(value any) error {
				return validateIfResolvable(true, nil)
			}))),
		validation.Field(&o.Populate, validation.Each(validation.By(
			func(value any) error {
				v := value.(orm_schema.Populate)
				return v.Validate()
			}))),
	)
}
