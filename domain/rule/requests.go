package rule

import (
	"ifttt/manager/domain/condition"
	"ifttt/manager/domain/resolvable"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetRulesRequest struct {
	Name string `mapstructure:"name" json:"name"`
}

func (g *GetRulesRequest) Validate() error {
	return validation.Validate(&g.Name)
}

type CreateRuleRequest struct {
	Name        string                  `json:"name" mapstructure:"name"`
	Description string                  `json:"description" mapstructure:"description"`
	Pre         []resolvable.Resolvable `json:"pre" mapstructure:"pre"`
	Switch      RuleSwitch              `json:"switch" mapstructure:"switch"`
}

func (c *CreateRuleRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Description, validation.Required, validation.Length(3, 0)),
		validation.Field(&c.Pre, validation.Each(validation.By(func(value interface{}) error {
			r := value.(resolvable.Resolvable)
			return r.Validate()
		}))),
		validation.Field(&c.Switch, validation.Required, validation.By(func(value interface{}) error {
			rs := value.(RuleSwitch)
			return rs.Validate()
		})),
	)
}

func (rs *RuleSwitch) Validate() error {
	return validation.ValidateStruct(rs,
		validation.Field(&rs.Cases, validation.Each(validation.By(func(value interface{}) error {
			rsw := value.(RuleSwitchCase)
			return rsw.Validate()
		}))),
		validation.Field(&rs.Default, validation.Required, validation.By(func(value interface{}) error {
			rsw := value.(RuleSwitchCase)
			return rsw.Validate()
		})),
	)
}

func (rsw *RuleSwitchCase) Validate() error {
	return validation.ValidateStruct(rsw,
		validation.Field(&rsw.Condition, validation.By(func(value interface{}) error {
			c := value.(condition.Condition)
			return c.Validate()
		})),
		validation.Field(&rsw.Do, validation.Each(validation.By(func(value interface{}) error {
			r := value.(resolvable.Resolvable)
			return r.Validate()
		}))),
		validation.Field(&rsw.Return, validation.By(func(value interface{}) error {
			r := value.(resolvable.Resolvable)
			return r.Validate()
		})),
	)
}
