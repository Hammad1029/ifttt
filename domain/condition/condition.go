package condition

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/resolvable"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Condition struct {
	ConditionType string                 `json:"conditionType" mapstructure:"conditionType"`
	Conditions    []Condition            `json:"conditions" mapstructure:"conditions"`
	Group         bool                   `json:"group" mapstructure:"group"`
	Operator1     *resolvable.Resolvable `json:"op1" mapstructure:"op1"`
	Operand       string                 `json:"opnd" mapstructure:"opnd"`
	Operator2     *resolvable.Resolvable `json:"op2" mapstructure:"op2"`
}

var (
	conditionTypeAnd = "AND"
	conditionTypeOr  = "OR"
	operandTypes     = []string{"eq", "ne", "in", "notIn", "lt", "lte", "gt", "gte"}
)

func (c *Condition) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ConditionType,
			validation.When(c.Group, validation.In(conditionTypeAnd, conditionTypeOr)).Else(validation.Empty)),
		validation.Field(&c.Conditions,
			validation.When(c.Group, validation.Length(1, 0),
				validation.Each(validation.By(func(value interface{}) error {
					c := value.(Condition)
					return c.Validate()
				}))).Else(validation.Nil)),
		validation.Field(&c.Group),
		validation.Field(&c.Operator1, validation.When(!c.Group, validation.By(
			func(value interface{}) error {
				r := value.(*resolvable.Resolvable)
				if r == nil {
					return validation.NewError("resolvable-not-nil", "Operator resolvable cannot be nil")
				}
				return r.Validate()
			}),
		).Else(validation.Nil)),

		validation.Field(&c.Operand, validation.When(!c.Group, validation.In(
			common.ConvertStringToInterfaceArray(operandTypes)...)).Else(validation.Empty)),
		validation.Field(&c.Operator2, validation.When(!c.Group, validation.By(
			func(value interface{}) error {
				r := value.(*resolvable.Resolvable)
				if r == nil {
					return validation.NewError("resolvable-not-nil", "Operator resolvable cannot be nil")
				}
				return r.Validate()
			}),
		).Else(validation.Nil)),
	)
}
