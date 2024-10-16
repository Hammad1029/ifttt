package rule

import (
	"ifttt/manager/domain/condition"
	"ifttt/manager/domain/resolvable"
)

type Rule struct {
	Id          uint                    `json:"id" mapstructure:"id"`
	Name        string                  `json:"name" mapstructure:"name"`
	Description string                  `json:"description" mapstructure:"description"`
	Pre         []resolvable.Resolvable `json:"pre" mapstructure:"pre"`
	Switch      RuleSwitch              `json:"switch" mapstructure:"switch"`
}

type RuleSwitch struct {
	Cases   []RuleSwitchCase `json:"cases" mapstructure:"cases"`
	Default RuleSwitchCase   `json:"default" mapstructure:"default"`
}

type RuleSwitchCase struct {
	Condition condition.Condition     `json:"condition" mapstructure:"condition"`
	Do        []resolvable.Resolvable `json:"do" mapstructure:"do"`
	Return    resolvable.Resolvable   `json:"return" mapstructure:"return"`
}
