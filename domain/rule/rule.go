package rule

import "ifttt/manager/domain/resolvable"

type Rule struct {
	Name        string                  `json:"name" mapstructure:"name"`
	Description string                  `json:"description" mapstructure:"description"`
	Conditions  Condition               `json:"conditions" mapstructure:"conditions"`
	Then        []resolvable.Resolvable `json:"then" mapstructure:"then"`
	Else        []resolvable.Resolvable `json:"else" mapstructure:"else"`
}

type Condition struct {
	ConditionType string                 `json:"conditionType" mapstructure:"conditionType"`
	Conditions    []Condition            `json:"conditions" mapstructure:"conditions"`
	Group         bool                   `json:"group" mapstructure:"group"`
	Operator1     *resolvable.Resolvable `json:"op1" mapstructure:"op1"`
	Operand       string                 `json:"opnd" mapstructure:"opnd"`
	Operator2     *resolvable.Resolvable `json:"op2" mapstructure:"op2"`
}
