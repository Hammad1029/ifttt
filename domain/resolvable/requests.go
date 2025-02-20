package resolvable

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/orm_schema"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var castError = validation.NewError("resolvable_not_casted", "could not cast resolvable")

func (r *Resolvable) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ResolveType, validation.Required,
			validation.In(common.ConvertStringToInterfaceArray(resolveTypes)...)),
		validation.Field(&r.ResolveData, validation.NotNil, validation.By(
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

func (u *uuid) Validate() error {
	return nil
}

func (g *getCache) Validate() error {
	return validation.Validate(&g.Key, validation.By(func(value any) error {
		return mustBeResolvable(value)
	}))
}

func (g *deleteCache) Validate() error {
	return validation.Validate(&g.Key, validation.By(func(value any) error {
		return mustBeResolvable(value)
	}))
}
func (s *setCache) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Key, validation.Required, validation.By(func(value any) error {
			r, ok := value.(Resolvable)
			if !ok {
				return castError
			}

			return r.Validate()
		})),
		validation.Field(&s.Value, validation.Required, validation.By(func(value any) error {
			r, ok := value.(Resolvable)
			if !ok {
				return castError
			}

			return r.Validate()
		})),
		validation.Field(&s.TTL),
	)
}

func (s *encode) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Alg, validation.Required,
			validation.In(common.EncodeMD5, common.EncodeSHA1, common.EncodeSHA2,
				common.EncodeBcrypt, common.EncodeBase64Decode, common.EncodeBase64Encode),
		),
		validation.Field(&s.Input, validation.Required, validation.By(func(value any) error {
			r, ok := value.(Resolvable)
			if !ok {
				return castError
			}

			return r.Validate()
		})))
}

func (s *setLog) Validate() error {
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

func (g *getErrors) Validate() error {
	return nil
}

func (g *getStore) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.Query, validation.Required, validation.By(func(value any) error {
			if _, ok := value.(string); ok {
				return nil
			}
			if r, ok := value.(map[string]any); ok {
				return mustBeResolvable(r)
			}
			return validation.NewError("jq_invalid_query", "jq query must be string or resolvable")
		})),
	)
}

func (g *getHeaders) Validate() error {
	return nil
}

func (g *getConst) Validate() error {
	return nil
}

func (c *apiCall) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Method, validation.Required, validation.In(
			http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut)),
		validation.Field(&c.URL, validation.Required, validation.By(
			func(value any) error {
				r, ok := value.(Resolvable)
				if !ok {
					return castError
				}

				return r.Validate()
			})),
		validation.Field(&c.Headers, validation.By(func(value any) error {
			m := value.(map[string]any)
			return ValidateIfResolvable(m)
		})),
		validation.Field(&c.Body, validation.By(func(value any) error {
			m := value.(map[string]any)
			return ValidateIfResolvable(m)
		})),
		validation.Field(&c.Aysnc),
		validation.Field(&c.Timeout),
	)
}

func (c *arithmetic) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Group, validation.NotNil),
		validation.Field(&c.Operation, validation.When(c.Group, validation.Required,
			validation.In(
				common.CalculatorAdd, common.CalculatorSubtract, common.CalculatorMultiply,
				common.CalculatorDivide, common.CalculatorModulus,
			)).Else(validation.Empty)),
		validation.Field(&c.Operators, validation.Each(validation.Required, validation.By(
			func(value any) error {
				a := value.(arithmetic)
				return a.Validate()
			}))),
		validation.Field(&c.Value, validation.When(!c.Group, validation.By(
			func(value any) error {
				r, ok := value.(*Resolvable)
				if !ok {
					return castError
				}

				return r.Validate()
			})).Else(validation.Nil)),
	)
}

func (c *jq) Validate() error {
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

func (c *stringInterpolation) Validate() error {
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
					return castError
				}

				return r.Validate()
			}))),
	)
}

func (c *query) Validate() error {
	count := len(common.RegexPositionalParameters.FindAllString(c.QueryString, -1))
	return validation.ValidateStruct(c,
		validation.Field(&c.QueryString, validation.Required),
		validation.Field(&c.Parameters,
			validation.Length(count, count), validation.Each(validation.By(
				func(value any) error {
					param := value.(Resolvable)
					return param.Validate()
				},
			))),
		validation.Field(&c.Async),
		validation.Field(&c.Timeout),
	)
}

func (e *response) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Event, validation.In(
			common.EventSuccess, common.EventExhaust, common.EventBadRequest, common.EventNotFound, common.EventSystemMalfunction,
		)),
	)
}

func (c *setStore) Validate() error {
	var mapCasted map[string]any = *c
	return ValidateIfResolvable(mapCasted)
}

func (c *cast) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.To, validation.In(common.ConvertStringToInterfaceArray(
			[]string{common.CastToString, common.CastToNumber, common.CastToBoolean})...)),
		validation.Field(&c.Input, validation.By(func(value any) error {
			return ValidateIfResolvable(value)
		})),
	)
}

func (o *Orm) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Query, validation.Nil),
		validation.Field(&o.SuccessiveQuery, validation.Nil),
		validation.Field(&o.Operation, is.UpperCase, validation.In(
			common.OrmSelect, common.OrmInsert, common.OrmUpdate, common.OrmDelete)),
		validation.Field(&o.Model, validation.NotNil),
		validation.Field(&o.Project,
			validation.When(o.Operation == common.OrmSelect,
				validation.By(
					func(value any) error {
						if casted, ok := value.(*[]orm_schema.Projection); !ok {
							return castError
						} else if casted != nil {
							for _, p := range *casted {
								if err := p.Validate(false); err != nil {
									return err
								}
							}
						}
						return nil
					}),
			).Else(validation.Nil)),
		validation.Field(&o.Columns,
			validation.When(o.Operation == common.OrmInsert || o.Operation == common.OrmUpdate,
				validation.Each(validation.By(func(value any) error {
					return ValidateIfResolvable(value)
				})),
			).Else(validation.Nil)),
		validation.Field(&o.Populate,
			validation.When(o.Operation == common.OrmSelect, validation.By(
				func(value any) error {
					if casted, ok := value.(*[]orm_schema.Populate); !ok {
						return castError
					} else if casted != nil {
						for _, p := range *casted {
							if err := p.Validate(ValidateIfResolvable); err != nil {
								return err
							}
						}
					}
					return nil
				}),
			).Else(validation.Nil)),
		validation.Field(&o.Where,
			validation.When(o.Operation == common.OrmSelect || o.Operation == common.OrmUpdate || o.Operation == common.OrmDelete,
				validation.By(
					func(value any) error {
						if v, ok := value.(*orm_schema.Where); !ok {
							return castError
						} else {
							return v.Validate(ValidateIfResolvable)
						}
					}),
			).Else(validation.Nil)),
		validation.Field(&o.OrderBy, validation.When(o.Operation != common.OrmSelect, validation.Empty)),
		validation.Field(&o.Limit, validation.When(o.Operation != common.OrmSelect, validation.Empty)),
	)
}

func (d *dateFunc) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Input,
			validation.By(func(value any) error {
				v := value.(dateInput)
				return v.Validate()
			})),
		validation.Field(&d.Manipulators, validation.Each(
			validation.By(func(value any) error {
				m := value.(dateManipulator)
				return m.Validate()
			}))),
		validation.Field(&d.Format),
		validation.Field(&d.UTC),
	)
}

func (d *dateManipulator) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Operator, validation.In(common.DateOperatorAdd, common.DateOperatorSubtract)),
		validation.Field(&d.Operand,
			validation.By(func(value any) error {
				v := value.(Resolvable)
				return v.Validate()
			})),
		validation.Field(&d.Unit, validation.Required, validation.In(
			common.ConvertStringToInterfaceArray(common.DateManipulatorUnits)...)),
	)
}

func (d *dateInput) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Input, validation.When(d.Input != nil, validation.By(
			func(value any) error {
				r := value.(*Resolvable)
				return r.Validate()
			}))),
		validation.Field(&d.Parse, validation.When(d.Input == nil, validation.Empty)),
		validation.Field(&d.Timezone),
	)
}

func (f *filterMap) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Input, validation.NotNil, validation.By(
			func(value any) error {
				return ValidateIfResolvable(value)
			})),
		validation.Field(&f.Do, validation.NotNil, validation.By(
			func(value any) error {
				return ValidateIfResolvable(value)
			})),
		validation.Field(&f.Condition, validation.Required, validation.By(func(value any) error {
			r := value.(Condition)
			return r.Validate()
		})),

		validation.Field(&f.Async),
	)
}

func (f *getIter) Validate() error {
	return nil
}

func (c *Condition) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ConditionType,
			validation.When(c.Group, validation.In(common.ConditionTypeAnd, common.ConditionTypeOr)).Else(validation.Empty)),
		validation.Field(&c.Conditions,
			validation.When(c.Group, validation.Length(1, 0),
				validation.Each(validation.By(func(value interface{}) error {
					c := value.(Condition)
					return c.Validate()
				}))).Else(validation.Nil)),
		validation.Field(&c.Group),
		validation.Field(&c.ComparisionType, validation.When(!c.Group, validation.Required,
			validation.In(
				common.ComparisionTypeString, common.ComparisionTypeNumber,
				common.ComparisionTypeBoolean, common.ComparisionTypeDate, common.ComparisionTypeBcrypt,
			)).Else(validation.Empty)),
		validation.Field(&c.Operator1, validation.When(!c.Group, validation.By(
			func(value interface{}) error {
				r := value.(*Resolvable)
				if r == nil {
					return validation.NewError("resolvable-not-nil", "Operator resolvable cannot be nil")
				}
				return r.Validate()
			}),
		).Else(validation.Nil)),
		validation.Field(&c.Operand, validation.When(!c.Group, validation.In(
			common.ConvertStringToInterfaceArray(common.OperandTypes)...)).Else(validation.Empty)),
		validation.Field(&c.Operator2, validation.When(!c.Group, validation.By(
			func(value interface{}) error {
				r := value.(*Resolvable)
				if r == nil {
					return validation.NewError("resolvable-not-nil", "Operator resolvable cannot be nil")
				}
				return r.Validate()
			}),
		).Else(validation.Nil)),
	)
}

func (c *conditional) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Condition, validation.Required, validation.By(func(value any) error {
			r := value.(Condition)
			return r.Validate()
		})),
		validation.Field(&c.True,
			validation.Required, validation.Each(validation.By(
				func(value any) error {
					param := value.(Resolvable)
					return param.Validate()
				},
			))),
		validation.Field(&c.False,
			validation.Required, validation.Each(validation.By(
				func(value any) error {
					param := value.(Resolvable)
					return param.Validate()
				},
			))),
	)
}

func (d *dateIntervals) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Start, validation.By(func(value any) error {
			input := value.(dateInput)
			return input.Validate()
		})),
		validation.Field(&d.End, validation.By(func(value any) error {
			input := value.(dateInput)
			return input.Validate()
		})),
		validation.Field(&d.Unit, validation.Required, validation.In(
			common.ConvertStringToInterfaceArray(common.DateManipulatorUnits)...)),
		validation.Field(&d.Format),
	)

}
