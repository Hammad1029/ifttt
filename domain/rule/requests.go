package rule

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/resolvable"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	conditionTypeAnd = "AND"
	conditionTypeOr  = "OR"
	operandTypes     = []string{"eq", "ne", "in", "notIn", "lt", "lte", "gt", "gte"}
)

type CreateRuleRequest struct {
	Name        string                  `json:"name" mapstructure:"name"`
	Description string                  `json:"description" mapstructure:"description"`
	Conditions  Condition               `json:"conditions" mapstructure:"conditions"`
	Then        []resolvable.Resolvable `json:"then" mapstructure:"then"`
	Else        []resolvable.Resolvable `json:"else" mapstructure:"else"`
}

func (c *CreateRuleRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Conditions, validation.Required, validation.By(func(value interface{}) error {
			c := value.(Condition)
			return c.Validate()
		})),
		validation.Field(&c.Then, validation.Required, validation.Each(validation.By(func(value interface{}) error {
			r := value.(resolvable.Resolvable)
			return r.Validate()
		}))),
		validation.Field(&c.Else, validation.Each(validation.By(func(value interface{}) error {
			r := value.(resolvable.Resolvable)
			return r.Validate()
		}))),
	)
}

func (c *Condition) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ConditionType,
			validation.When(c.Group, validation.In(conditionTypeAnd, conditionTypeOr)).Else(validation.Empty)),
		validation.Field(&c.Conditions,
			validation.When(c.Group, validation.Length(1, 0),
				validation.Each(validation.By(func(value interface{}) error {
					c := value.(Condition)
					return c.Validate()
				}))).Else(validation.Empty)),
		validation.Field(&c.Group),
		validation.Field(&c.Operator1, validation.When(!c.Group, validation.Length(1, 0),
			validation.Each(validation.By(func(value interface{}) error {
				r := value.(resolvable.Resolvable)
				return r.Validate()
			}))).Else(validation.Empty)),
		validation.Field(&c.Operand, validation.When(!c.Group, validation.In(
			common.ConvertStringToInterfaceArray(operandTypes)...)).Else(validation.Empty)),
		validation.Field(&c.Operator2, validation.When(!c.Group, validation.Length(1, 0),
			validation.Each(validation.By(func(value interface{}) error {
				r := value.(resolvable.Resolvable)
				return r.Validate()
			}))).Else(validation.Empty)),
	)
}
